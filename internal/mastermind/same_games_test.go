package mastermind_test

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
)

// everyCode offers every code not played yet. It is what GetActions must give the same games as.
func everyCode(gameplay game.GamePlay, remainingWorlds game.PossibleWorlds, aa game.AvailableActions) []game.Action {
	actions := []game.Action{}
	for _, code := range aa.WorldActions() {
		if !gameplay.PreviouslyActed(code) {
			actions = append(actions, code)
		}
	}
	return actions
}

var everyCodeAgreement = game.Agreement{GetActions: everyCode, Act: mastermind.Act}

func clearCaches() {
	game.ThePredictionsFreeCache.Clear()
	game.ObserveCache.Cache.Clear()
}

// position plays the given guesses and responses from the start and returns the game so far and the
// codes still possible.
func position(lines []string, turns ...string) (game.GamePlay, game.PossibleWorlds) {
	game.BUniverse = game.InitUniverse(lines)
	remaining := game.BUniverse.PossibleWorlds
	gameplay := game.InitGamePlay(remaining.GetHash())
	for _, turn := range turns {
		guessAndResponse := strings.Split(turn, ":")
		reaction := game.NewReaction(game.NewAction(guessAndResponse[0]), guessAndResponse[1], 0.0)
		gameplay, remaining = play.Observe(gameplay, reaction, mastermind.Act, remaining)
	}
	return gameplay, remaining
}

// split counts how many of the remaining codes give each response to guess.
func split(gameplay game.GamePlay, remaining game.PossibleWorlds, guess game.Action) string {
	counts := map[string]int{}
	for _, index := range remaining.Indices() {
		counts[mastermind.Act(gameplay, game.BUniverse.GetWorldByIndex(index), guess).Response.String()]++
	}
	parts := []string{}
	for response, count := range counts {
		parts = append(parts, fmt.Sprintf("%s:%d", response, count))
	}
	sort.Strings(parts)
	return strings.Join(parts, " ")
}

func offered(actions []game.Action) map[string]bool {
	isOffered := map[string]bool{}
	for _, action := range actions {
		isOffered[action.String()] = true
	}
	return isOffered
}

// After ABCD scores one black and ABCC scores nothing, 27 codes remain. AEFE and AEEF split them into
// different groups, so both have to be offered.
func TestGetActionsKeepsGuessesThatSplitDifferently(t *testing.T) {
	lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")
	clearCaches()
	gameplay, remaining := position(lines, "ABCD:B", "ABCC:!")
	if remaining.Len() != 27 {
		t.Fatalf("expected 27 codes to remain, got %d", remaining.Len())
	}
	if split(gameplay, remaining, game.NewAction("AEFE")) == split(gameplay, remaining, game.NewAction("AEEF")) {
		t.Fatalf("expected AEFE and AEEF to split the remaining codes differently")
	}
	isOffered := offered(mastermind.GetActions(gameplay, remaining, game.BUniverse.Actions))
	for _, code := range []string{"AEEF", "AEFE"} {
		if !isOffered[code] {
			t.Errorf("%s is not offered", code)
		}
	}
	for _, code := range everyCode(gameplay, remaining, game.BUniverse.Actions) {
		if isOffered[code.String()] {
			continue
		}
		found := false
		for kept := range isOffered {
			if split(gameplay, remaining, game.NewAction(kept)) == split(gameplay, remaining, code) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%s is not offered and no offered code splits the remaining codes as it does", code)
		}
	}
}

type fixture struct {
	name  string
	lines func() []string
	heavy bool
}

func fromFile(path string) func() []string {
	return func() []string { return utils.ReadLinesFromFilePath(path) }
}

// someOf returns count codes of standard Mastermind picked at random, always the same ones, in their
// order in the file.
func someOf(count int, seed int64) func() []string {
	return func() []string {
		lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")
		picked := rand.New(rand.NewSource(seed)).Perm(len(lines))[:count]
		sort.Ints(picked)
		some := []string{}
		for _, i := range picked {
			some = append(some, lines[i])
		}
		return some
	}
}

// noRepeatedColour returns the 360 codes of standard Mastermind that use four different colours. The
// list is not every combination of the colours, but moving pegs and renaming colours turns it into
// itself, so GetActions can still leave codes out.
func noRepeatedColour() []string {
	some := []string{}
	for _, line := range utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt") {
		if strings.Count(line, line[:1]) == 1 && strings.Count(line, line[1:2]) == 1 && strings.Count(line, line[2:3]) == 1 && strings.Count(line, line[3:4]) == 1 {
			some = append(some, line)
		}
	}
	return some
}

var fixtures = []fixture{
	{"3 pegs 4 colours", fromFile("../../fixtures/mastermind-ABCD-3.txt"), false},
	{"100 random codes of 4 pegs 6 colours", someOf(100, 1), false},
	{"300 random codes of 4 pegs 6 colours", someOf(300, 2), false},
	{"4 pegs 6 colours without a repeated colour", noRepeatedColour, true},
	{"4 pegs 6 colours", fromFile("../../fixtures/commercial-mastermind.txt"), true},
	{"4 pegs 7 colours", fromFile("../../fixtures/mastermind-ABCDEFG-4.txt"), true},
}

var weightings = []struct {
	name       string
	motivation game.Motivation
}{
	{"MostParts", game.Motivation{MostParts: 1.0}},
	{"Entropy", game.Motivation{Entropy: 1.0}},
	{"MiniMax", game.Motivation{MiniMax: 1.0}},
	{"Reward", game.Motivation{Reward: 1.0}},
	{"Reward 0.75 Entropy 0.25", game.Motivation{Reward: 0.75, Entropy: 0.25}},
	{"Entropy 0.05 MiniMax 0.70 MostParts 0.25", game.Motivation{Entropy: 0.05, MiniMax: 0.70, MostParts: 0.25}},
	{"Entropy 0.85 MostParts 0.15", game.Motivation{Entropy: 0.85, MostParts: 0.15}},
}

func turns(gameplay game.GamePlay) string {
	played := []string{}
	for _, reaction := range gameplay.Reactions {
		played = append(played, reaction.Action.String()+":"+reaction.Response.String())
	}
	return strings.Join(played, " ")
}

// Every game must be played exactly as it is when every code not played yet is offered.
func TestGetActionsPlaysTheSameGamesAsEveryCode(t *testing.T) {
	for _, f := range fixtures {
		for _, w := range weightings {
			f, w := f, w
			t.Run(f.name+"/"+w.name, func(t *testing.T) {
				if f.heavy && testing.Short() {
					t.Skip("plays every game twice; skipped with -short")
				}
				lines := f.lines()
				clearCaches()
				expected := play.PlayGames(everyCodeAgreement, w.motivation, lines, len(lines))
				clearCaches()
				got := play.PlayGames(mastermind.Agreement, w.motivation, lines, len(lines))
				different, gotTotal, expectedTotal, gotLongest, expectedLongest := 0, 0, 0, 0, 0
				for i := range lines {
					gotTotal += got[i].TurnCount
					expectedTotal += expected[i].TurnCount
					if got[i].TurnCount > gotLongest {
						gotLongest = got[i].TurnCount
					}
					if expected[i].TurnCount > expectedLongest {
						expectedLongest = expected[i].TurnCount
					}
					if turns(got[i]) != turns(expected[i]) {
						different++
						if different <= 3 {
							t.Errorf("answer %s: played %s, with every code offered %s", lines[i], turns(got[i]), turns(expected[i]))
						}
					}
				}
				t.Logf("%d guesses, longest game %d; with every code offered %d guesses, longest game %d", gotTotal, gotLongest, expectedTotal, expectedLongest)
				if different > 0 {
					t.Errorf("%d of %d games differ", different, len(lines))
				}
			})
		}
	}
}

// At every position the games reach, each code that is not offered must split the remaining codes
// exactly as an offered code that comes before it alphabetically, with exactly the same scores, so
// Decide never prefers the code that was left out.
func TestGetActionsLeavesOutOnlyCodesThatScoreLikeAnEarlierOfferedCode(t *testing.T) {
	for _, f := range fixtures {
		f := f
		t.Run(f.name, func(t *testing.T) {
			if f.heavy && testing.Short() {
				t.Skip("scores every code at every position; skipped with -short")
			}
			lines := f.lines()
			type reached struct {
				turns     []string
				gameplay  game.GamePlay
				remaining game.PossibleWorlds
			}
			positions := map[uint64]reached{}
			for _, w := range weightings {
				clearCaches()
				for _, g := range play.PlayGames(everyCodeAgreement, w.motivation, lines, len(lines)) {
					remaining := game.BUniverse.PossibleWorlds
					gameplay := game.InitGamePlay(remaining.GetHash())
					turnStrings := []string{}
					for _, reaction := range g.Reactions {
						positions[gameplay.Hash] = reached{append([]string{}, turnStrings...), gameplay, remaining}
						turnStrings = append(turnStrings, reaction.Action.String()+":"+reaction.Response.String())
						gameplay, remaining = play.Observe(gameplay, reaction, mastermind.Act, remaining)
					}
				}
			}
			checked, left := 0, 0
			for _, p := range positions {
				gameplay, remaining, turnStrings := p.gameplay, p.remaining, p.turns
				if remaining.Len() < 2 {
					continue
				}
				all := everyCode(gameplay, remaining, game.BUniverse.Actions)
				offeredCodes := mastermind.GetActions(gameplay, remaining, game.BUniverse.Actions)
				checked++
				if len(offeredCodes) == len(all) {
					continue
				}
				scores := game.ExploreReactions(gameplay, remaining, all, mastermind.Act)
				isOffered := offered(offeredCodes)
				byScore := map[string][]string{}
				for _, code := range offeredCodes {
					key := scoreKey(scores.Get(code.String())) + " " + split(gameplay, remaining, code)
					byScore[key] = append(byScore[key], code.String())
				}
				for _, code := range all {
					if isOffered[code.String()] {
						continue
					}
					left++
					key := scoreKey(scores.Get(code.String())) + " " + split(gameplay, remaining, code)
					earlier := false
					for _, kept := range byScore[key] {
						if kept < code.String() {
							earlier = true
						}
					}
					if !earlier {
						t.Fatalf("after %v: %s is left out, but no offered code before it scores the same (%s)", turnStrings, code, key)
					}
				}
			}
			t.Logf("%d positions checked, %d codes left out across them", checked, left)
		})
	}
}

func scoreKey(p game.Prediction) string {
	return fmt.Sprintf("reward %x entropy %x largest %x parts %x possible %v", p.Reward, p.Entropy, p.InversedMax, p.PartsCount, p.ActionMatchesAPossibleWorld)
}
