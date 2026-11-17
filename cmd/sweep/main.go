// cmd/sweep plays a list of weightings and prints one line of JSON per weighting, in the same shape as
// cmd/optimize. You either list the weightings with -motivations, or give -trial.stepSize to try every
// combination at that step, in the same order as cmd/optimize. Every game is played with the engine's
// own functions for scoring the guesses, picking one and narrowing the answers; nothing in the engine
// changes.
//
// -mode games (the default) plays the games one after another with play.PlayGames, as cmd/player does.
// -mode tree plays all the games of a weighting together. At each position (the guesses so far and
// their responses) it scores the guesses once with the engine's ExploreReactions, picks one with
// Decide, then splits the answers still possible by the response they give and carries on with each
// group. A position that many games reach is scored once instead of once per game, and its scores are
// never written to and read back from the cache. The guess chosen depends only on the position, so
// every game gets the same guesses as in -mode games; sweep_test.go checks that the totals agree.
// -tree.cacheStates keeps the scores of positions near the start (at most -tree.cacheDepth turns in)
// from one weighting to the next, since neighbouring weightings often reach the same early positions.
//
// -gamelog writes one JSON line per game. -shard.index and -shard.count split the weightings across
// separate processes.
package main

import (
	"bufio"
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/coocood/freecache"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/chrisconley/go-hangman/internal/wordle"
)

type trialResult struct {
	Motivation game.Motivation
	GameStats  play.Stats
	Seconds    float64 `json:",omitempty"`
}

type gameRecord struct {
	Trial     string   `json:"trial"`
	Word      string   `json:"word"`
	Turns     int      `json:"turns"`
	Pain      float64  `json:"pain"`
	Actions   []string `json:"actions"`
	Responses []string `json:"responses"`
}

func key(m game.Motivation) string {
	return fmt.Sprintf("%.3f|%.3f|%.3f|%.3f", m.Reward, m.Entropy, m.MiniMax, m.MostParts)
}

type stateCache struct {
	capacity int
	order    *list.List
	entries  map[uint64]*list.Element
}

type cacheEntry struct {
	hash        uint64
	predictions game.Predictions
}

func (c *stateCache) get(hash uint64) (game.Predictions, bool) {
	if element, ok := c.entries[hash]; ok {
		c.order.MoveToFront(element)
		return element.Value.(cacheEntry).predictions, true
	}
	return game.Predictions{}, false
}

func (c *stateCache) put(hash uint64, predictions game.Predictions) {
	if c.capacity <= 0 {
		return
	}
	c.entries[hash] = c.order.PushFront(cacheEntry{hash, predictions})
	if c.order.Len() > c.capacity {
		oldest := c.order.Back()
		c.order.Remove(oldest)
		delete(c.entries, oldest.Value.(cacheEntry).hash)
	}
}

type treePlayer struct {
	agreement  game.Agreement
	motivation game.Motivation
	cache      *stateCache
	cacheDepth int
	games      []game.GamePlay
	fast       *fastWordle
	check      bool

	checkedStates int
}

// branch gives the next position its own copy of the guesses and responses so far. GamePlay.AddTurn
// appends to Reactions, so positions that follow from the same earlier position would otherwise share,
// and overwrite, one backing array. play.PlayGames never meets this because each game has its own
// GamePlay.
func branch(gameplay game.GamePlay) game.GamePlay {
	gameplay.Reactions = append(make([]game.Reaction, 0, len(gameplay.Reactions)+1), gameplay.Reactions...)
	return gameplay
}

func (t *treePlayer) explore(gameplay game.GamePlay, possibleWorlds game.PossibleWorlds) game.Predictions {
	cacheable := gameplay.TurnCount <= t.cacheDepth
	if cacheable {
		if predictions, ok := t.cache.get(gameplay.Hash); ok {
			return predictions
		}
	}
	actions := t.agreement.GetActions(gameplay, possibleWorlds, game.BUniverse.Actions)
	predictions := game.ExploreReactions(gameplay, possibleWorlds, actions, t.agreement.Act)
	if cacheable {
		t.cache.put(gameplay.Hash, predictions)
	}
	return predictions
}

func (t *treePlayer) play(gameplay game.GamePlay, possibleWorlds game.PossibleWorlds) {
	predictions := t.explore(gameplay, possibleWorlds)
	action := game.Decide(predictions, t.motivation).Action
	groups := map[string][]uint64{}
	reactions := map[string]game.Reaction{}
	order := []string{}
	for _, worldIndex := range possibleWorlds.Indices() {
		reaction := t.agreement.Act(gameplay, game.BUniverse.GetWorldByIndex(worldIndex), action)
		response := reaction.Response.String()
		if first, ok := reactions[response]; !ok {
			reactions[response] = reaction
			order = append(order, response)
		} else if first.Reward != reaction.Reward {
			panic("two answers gave the same response but different rewards; -mode tree needs the reward to follow from the response")
		}
		groups[response] = append(groups[response], worldIndex)
	}
	for _, response := range order {
		nextGameplay, nextWorlds := play.Observe(branch(gameplay), reactions[response], t.agreement.Act, possibleWorlds)
		if nextWorlds.Len() != len(groups[response]) {
			panic("Observe and the split by response disagree on which answers are left")
		}
		if len(groups[response]) == 1 && action.Matches(game.BUniverse.GetWorldByIndex(groups[response][0])) {
			t.games[groups[response][0]] = nextGameplay
			continue
		}
		t.play(nextGameplay, nextWorlds)
	}
}

// trialSettings holds the flags that decide how one weighting's games are played.
type trialSettings struct {
	mode       string
	limit      int
	cacheSize  int
	cache      *stateCache
	cacheDepth int
	fast       *fastWordle
	check      bool
}

// playTrial plays the games of one weighting and returns them in input order, with the number of
// positions -tree.check compared against the engine.
func playTrial(agreement game.Agreement, lines []string, motivation game.Motivation, s trialSettings) ([]game.GamePlay, int) {
	game.ObserveCache.Cache = freecache.NewCache(500 * 1024 * 1024)
	if s.mode != "tree" {
		game.ThePredictionsFreeCache = freecache.NewCache(s.cacheSize * 1024 * 1024)
		return play.PlayGames(agreement, motivation, lines, s.limit), 0
	}
	game.BUniverse = game.InitUniverse(lines)
	player := &treePlayer{agreement: agreement, motivation: motivation, cache: s.cache, cacheDepth: s.cacheDepth, games: make([]game.GamePlay, len(lines)), fast: s.fast, check: s.check}
	root := game.InitGamePlay(game.BUniverse.PossibleWorlds.GetHash())
	if s.fast != nil {
		player.playFast(root, game.BUniverse.PossibleWorlds, make([]bool, len(s.fast.guesses)))
	} else {
		player.play(root, game.BUniverse.PossibleWorlds)
	}
	return player.games, player.checkedStates
}

// buildTrials returns every combination of weights at stepSize (none when it is 0), followed by the
// reward:entropy:minimax:mostparts weightings listed in motivations.
func buildTrials(stepSize float64, motivations string, mostPartsZero bool) []game.Motivation {
	trials := []game.Motivation{}
	if stepSize > 0 {
		for _, trial := range play.BuildTrials(0.0, 1.0, stepSize, 4) {
			for index, value := range trial {
				if value < 0.0 {
					trial[index] = 0.0
				}
			}
			m := game.Motivation{Reward: trial[3], Entropy: trial[0], MiniMax: trial[1], MostParts: trial[2]}
			if mostPartsZero && m.MostParts != 0 {
				continue
			}
			trials = append(trials, m)
		}
	}
	for _, spec := range strings.Split(motivations, ",") {
		if spec == "" {
			continue
		}
		parts := strings.Split(spec, ":")
		if len(parts) != 4 {
			panic("motivation tuple must be reward:entropy:minimax:mostparts")
		}
		values := make([]float64, 4)
		for i, p := range parts {
			v, err := strconv.ParseFloat(p, 64)
			if err != nil {
				panic(err)
			}
			values[i] = v
		}
		trials = append(trials, game.Motivation{Reward: values[0], Entropy: values[1], MiniMax: values[2], MostParts: values[3]})
	}
	return trials
}

// shard returns the weightings that process index of count should play: every count-th one, or one
// block of neighbours when contiguous is set.
func shard(trials []game.Motivation, index int, count int, contiguous bool) []game.Motivation {
	if count <= 1 {
		return trials
	}
	sharded := []game.Motivation{}
	for i, m := range trials {
		block := i * count / len(trials)
		if (contiguous && block == index) || (!contiguous && i%count == index) {
			sharded = append(sharded, m)
		}
	}
	return sharded
}

func main() {
	gameString := flag.String("game", "hangman", "hangman, mastermind or wordle")
	wordleGuesses := flag.String("wordle.guesses", "", "with -game wordle: file of allowed guesses, one per line (the possible answers come on stdin)")
	mode := flag.String("mode", "games", "games (one game after another, as cmd/player) or tree (all games of a weighting together, each position scored once)")
	cacheStates := flag.Int("tree.cacheStates", 0, "with -mode tree: keep the scores of up to this many early positions from one weighting to the next, dropping the least recently used")
	cacheDepth := flag.Int("tree.cacheDepth", 1, "with -mode tree: keep only positions reached after at most this many turns")
	fast := flag.Bool("tree.fast", false, "with -mode tree -game wordle: score guesses from a table of responses worked out in advance; the results are exactly the same numbers (see fast_wordle.go)")
	check := flag.Bool("tree.check", false, "with -tree.fast: also score and choose with the engine at every position, and stop with an error on any difference")
	stepSize := flag.Float64("trial.stepSize", 0, "try every combination of the four weights at this step, as cmd/optimize does (0 for none)")
	motivations := flag.String("motivations", "", "weightings to play, comma-separated, each as reward:entropy:minimax:mostparts, e.g. 0.75:0.25:0:0,1:0:0:0")
	mostPartsZero := flag.Bool("mostparts.zero", false, "with -trial.stepSize: keep only the combinations whose MostParts weight is 0")
	gameLog := flag.String("gamelog", "", "write each game's guesses and responses to this file, one JSON line per game")
	cacheSize := flag.Int("cacheSize", 1000, "with -mode games: size of the cache of scored positions in MB (the cmd/player default)")
	limit := flag.Int("limit", 0, "games per weighting; 0 plays every word on stdin")
	shardIndex := flag.Int("shard.index", 0, "with -shard.count: which share of the weightings this process plays, from 0")
	shardCount := flag.Int("shard.count", 1, "split the weightings across this many separate processes")
	shardContiguous := flag.Bool("shard.contiguous", false, "give each process a block of neighbouring weightings instead of every shard.count-th one")
	flag.Parse()

	lines, err := utils.ReadLinesFromStdin()
	if err != nil || len(lines) == 0 {
		fmt.Fprintln(os.Stderr, "expected one word per line on stdin")
		os.Exit(2)
	}
	if *limit <= 0 || *limit > len(lines) {
		*limit = len(lines)
	}
	agreement := hangman.Agreement
	if *gameString == "mastermind" {
		agreement = mastermind.Agreement
	} else if *gameString == "wordle" {
		if *wordleGuesses == "" {
			panic("-game wordle needs -wordle.guesses")
		}
		agreement = wordle.NewAgreement(utils.ReadLinesFromFilePath(*wordleGuesses))
	} else if *gameString != "hangman" {
		panic("game must be hangman, mastermind or wordle")
	}
	if *mode == "tree" && *limit != len(lines) {
		panic("-mode tree plays every answer on stdin; drop -limit or shorten the list")
	}
	debug.SetGCPercent(20)

	trials := shard(buildTrials(*stepSize, *motivations, *mostPartsZero), *shardIndex, *shardCount, *shardContiguous)
	fmt.Fprintf(os.Stderr, "words: %d, games per trial: %d, trials: %d, mode: %s\n", len(lines), *limit, len(trials), *mode)

	var logWriter *bufio.Writer
	if *gameLog != "" {
		f, err := os.Create(*gameLog)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		logWriter = bufio.NewWriterSize(f, 1<<20)
		defer logWriter.Flush()
	}

	cache := &stateCache{capacity: *cacheStates, order: list.New(), entries: map[uint64]*list.Element{}}
	var fastExplorer *fastWordle
	if *fast {
		if *gameString != "wordle" || *mode != "tree" {
			panic("-tree.fast needs -game wordle -mode tree")
		}
		fastExplorer = newFastWordle(utils.ReadLinesFromFilePath(*wordleGuesses), game.InitUniverse(lines).Worlds)
	}
	settings := trialSettings{mode: *mode, limit: *limit, cacheSize: *cacheSize, cache: cache, cacheDepth: *cacheDepth, fast: fastExplorer, check: *check}
	for _, motivation := range trials {
		start := time.Now()
		games, checkedStates := playTrial(agreement, lines, motivation, settings)
		if *check {
			fmt.Fprintf(os.Stderr, "%s: %d states checked against the engine's ExploreReactions and Decide\n", key(motivation), checkedStates)
		}
		stats := play.CalculateStats(games, true)
		js, err := json.Marshal(trialResult{Motivation: motivation, GameStats: stats, Seconds: time.Since(start).Seconds()})
		if err != nil {
			panic(err)
		}
		fmt.Println(string(js))
		if logWriter != nil {
			k := key(motivation)
			for j, gameplay := range games {
				rec := gameRecord{Trial: k, Word: game.BUniverse.GetWorldByIndex(uint64(j)).String(), Turns: gameplay.TurnCount, Pain: gameplay.TotalPain}
				for _, reaction := range gameplay.Reactions {
					rec.Actions = append(rec.Actions, reaction.Action.String())
					rec.Responses = append(rec.Responses, reaction.Response.String())
				}
				b, _ := json.Marshal(rec)
				logWriter.Write(b)
				logWriter.WriteString("\n")
			}
			logWriter.Flush()
		}
	}
}
