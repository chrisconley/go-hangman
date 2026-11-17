package main

// The -tree.fast explorer for Wordle. It computes, for every state the tree walk reaches, the same four
// statistics per action that game.ExploreReactions computes (Reward = summed per-world reward / worlds,
// Entropy = mth.CalcEntropy over the response-class counts in first-appearance order, InversedMax =
// mth.FindInversedMax, PartsCount = number of classes, valid = the guess is a remaining world), from a
// precomputed guess x answer pattern table instead of strings, maps and bitmaps, adding the per-world
// rewards in the same order so every float is bit-identical. It then picks the action with the same rule
// as game.Decide (argmax of Reward^a * Entropy^b * InversedMax^c * PartsCount^d multiplied left to
// right, ties to a valid action, then to the lexically smallest). -tree.check runs the engine's own
// ExploreReactions and Decide at every state as well and panics on any difference in a statistic or in
// the decision; sweep_test.go runs it.

import (
	"fmt"
	"math"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/mth"
	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/wordle"
)

const allGreen = 242

type fastWordle struct {
	guesses  []string
	index    map[string]int
	pattern  [][]uint8
	rewardOf [243]float64
	stats    []fastStat
	counts   [243]int
	order    [243]uint8
	buffer   [243]float64
}

type fastStat struct {
	reward, entropy, inversedMax, partsCount float64
	valid                                    bool
}

func newFastWordle(guesses []string, answers []game.World) *fastWordle {
	f := &fastWordle{guesses: guesses, index: map[string]int{}, stats: make([]fastStat, len(guesses))}
	for g, guess := range guesses {
		f.index[guess] = g
		row := make([]uint8, len(answers))
		for w, answer := range answers {
			code, _, _ := wordle.Score(answer.String(), guess)
			row[w] = uint8(code)
		}
		f.pattern = append(f.pattern, row)
	}
	for code := 0; code < 243; code++ {
		greens, yellows := 0, 0
		for c := code; c > 0; c /= 3 {
			if c%3 == 2 {
				greens++
			} else if c%3 == 1 {
				yellows++
			}
		}
		f.rewardOf[code] = wordle.Reward(greens, yellows)
	}
	return f
}

func (f *fastWordle) explore(acted []bool, worlds []uint64) {
	total := float64(len(worlds))
	for g, row := range f.pattern {
		if acted[g] {
			continue
		}
		parts := 0
		rewards := 0.0
		for _, w := range worlds {
			code := row[w]
			if f.counts[code] == 0 {
				f.order[parts] = code
				parts++
			}
			f.counts[code]++
			rewards += f.rewardOf[code]
		}
		for i := 0; i < parts; i++ {
			f.buffer[i] = float64(f.counts[f.order[i]])
		}
		f.stats[g] = fastStat{
			reward:      rewards / total,
			entropy:     mth.CalcEntropy(total, f.buffer[:parts]),
			inversedMax: mth.FindInversedMax(total, f.buffer[:parts]),
			partsCount:  float64(parts),
			valid:       f.counts[allGreen] > 0,
		}
		for i := 0; i < parts; i++ {
			f.counts[f.order[i]] = 0
		}
	}
}

func power(x float64, y float64) float64 {
	if y == 0 {
		return 1
	}
	return math.Pow(x, y)
}

func (f *fastWordle) decide(acted []bool, m game.Motivation) int {
	maxWeight := -1.0
	maxAction := -1
	maxIsValid := false
	for g := range f.guesses {
		if acted[g] {
			continue
		}
		s := f.stats[g]
		weight := power(s.reward, m.Reward) * power(s.entropy, m.Entropy) * power(s.inversedMax, m.MiniMax) * power(s.partsCount, m.MostParts)
		if weight == maxWeight && s.valid == maxIsValid && (maxAction < 0 || f.guesses[g] < f.guesses[maxAction]) {
			maxAction = g
		} else if weight == maxWeight && s.valid && !maxIsValid {
			maxAction = g
			maxIsValid = true
		} else if weight > maxWeight {
			maxWeight = weight
			maxAction = g
			maxIsValid = s.valid
		}
	}
	return maxAction
}

func (t *treePlayer) playFast(gameplay game.GamePlay, possibleWorlds game.PossibleWorlds, acted []bool) {
	f := t.fast
	worlds := possibleWorlds.Indices()
	f.explore(acted, worlds)
	g := f.decide(acted, t.motivation)
	if t.check {
		t.checkState(gameplay, possibleWorlds, acted, g)
	}
	action := game.NewAction(f.guesses[g])
	groups := map[uint8][]uint64{}
	order := []uint8{}
	for _, w := range worlds {
		code := f.pattern[g][w]
		if _, ok := groups[code]; !ok {
			order = append(order, code)
		}
		groups[code] = append(groups[code], w)
	}
	acted[g] = true
	for _, code := range order {
		reaction := t.agreement.Act(gameplay, game.BUniverse.GetWorldByIndex(groups[code][0]), action)
		nextGameplay, nextWorlds := play.Observe(branch(gameplay), reaction, t.agreement.Act, possibleWorlds)
		if nextWorlds.Len() != len(groups[code]) {
			panic("Observe disagrees with the pattern table")
		}
		if code == allGreen {
			t.games[groups[code][0]] = nextGameplay
			continue
		}
		t.playFast(nextGameplay, nextWorlds, acted)
	}
	acted[g] = false
}

func (t *treePlayer) checkState(gameplay game.GamePlay, possibleWorlds game.PossibleWorlds, acted []bool, g int) {
	f := t.fast
	actions := t.agreement.GetActions(gameplay, possibleWorlds, game.BUniverse.Actions)
	predictions := game.ExploreReactions(gameplay, possibleWorlds, actions, t.agreement.Act)
	available := 0
	for i := range f.guesses {
		if !acted[i] {
			available++
		}
	}
	if len(actions) != available {
		panic(fmt.Sprintf("GetActions offers %d actions, the fast explorer %d", len(actions), available))
	}
	for _, prediction := range predictions.PredictionsMap.Values {
		s := f.stats[f.index[prediction.Action.String()]]
		if s.reward != prediction.Reward || s.entropy != prediction.Entropy || s.inversedMax != prediction.InversedMax || s.partsCount != prediction.PartsCount || s.valid != prediction.ActionMatchesAPossibleWorld {
			panic(fmt.Sprintf("statistics differ for %s: fast %+v engine %+v", prediction.Action.String(), s, prediction))
		}
	}
	engine := game.Decide(predictions, t.motivation).Action.String()
	if engine != f.guesses[g] {
		panic(fmt.Sprintf("decision differs: fast %s engine %s", f.guesses[g], engine))
	}
	t.checkedStates++
}
