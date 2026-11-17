package wordle

import (
	"math"
	"os"
	"testing"

	"github.com/coocood/freecache"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/mth"
	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
)

// Each case was worked by hand: greens first, then yellows left to right while the answer still has an
// unmatched copy of the letter.
var feedbackCases = []struct{ answer, guess, want string }{
	{"abbey", "babes", "YYGG-"}, // b and a swapped, second b green, s absent
	{"sassy", "assay", "YYG-G"}, // three s in the answer: one green, one yellow; the second a has no copy left
	{"crane", "crane", "GGGGG"},
	{"cigar", "soare", "--YY-"},
	{"rebut", "soare", "---YY"},
	{"sissy", "soare", "G----"},
	{"hello", "llama", "YY---"}, // the answer's two l's are both unmatched, so both guessed l's are yellow
	{"hello", "lever", "YG---"}, // the answer's only e is green, so the second e is gray
	{"speed", "eerie", "YY---"}, // two e's in the answer, three in the guess: the third is gray
	{"abide", "speed", "--Y-Y"}, // one e in the answer: first e yellow, second gray
	{"erase", "speed", "Y-YY-"}, // two e's in the answer, both guessed e's yellow
	{"steal", "speed", "G-G--"}, // the green e uses the answer's only e, the other e is gray
	{"crepe", "speed", "-YGY-"}, // green e at 3, the answer's other e makes the e at 4 yellow
	{"elder", "eerie", "GYY--"},
	{"array", "rarer", "YYG--"}, // the green r uses one r, the first r takes the other, the last r is gray
	{"mamma", "magma", "GG-GG"},
	{"geese", "eerie", "YG--G"}, // greens at 2 and 5, the answer's third e makes the first e yellow, the r and i are gray
	{"otter", "torte", "YYYYY"},
	{"fuzzy", "pizza", "--GG-"},
	{"aback", "aaaaa", "G-G--"}, // two greens use both a's, the other three a's are gray
	{"skill", "lilly", "YY-G-"}, // the green l at 4 and the yellow l at 1 use both l's, the l at 3 is gray
	{"creep", "eerie", "YYY--"},
	{"those", "geese", "---GG"}, // greens are assigned first, so the early e's are gray, not yellow
	{"knoll", "lolly", "YY-G-"},
	{"tweet", "tenet", "GY-GG"}, // the green e at 4 leaves the answer's e at 3 for the guessed e at 2
}

func TestFeedback(t *testing.T) {
	for _, c := range feedbackCases {
		if got := Feedback(c.answer, c.guess); got != c.want {
			t.Errorf("Feedback(answer %q, guess %q) = %q, want %q", c.answer, c.guess, got, c.want)
		}
	}
}

func TestRewardAndAct(t *testing.T) {
	gp := game.InitGamePlay(0)
	reaction := Act(gp, game.NewWorld("abbey"), game.NewAction("babes"))
	if reaction.Response.String() != "YYGG-" || math.Abs(reaction.Reward-0.6) > 1e-12 {
		t.Errorf("abbey/babes: got %q reward %v, want YYGG- reward 0.6", reaction.Response.String(), reaction.Reward)
	}
	if r := Act(gp, game.NewWorld("crane"), game.NewAction("crane")).Reward; r != 1.0 {
		t.Errorf("a solving guess should be worth 1.0, got %v", r)
	}
	if r := Act(gp, game.NewWorld("fuzzy"), game.NewAction("crane")).Reward; r != 0.0 {
		t.Errorf("an all-gray guess should be worth 0, got %v", r)
	}
	if r := Reward(0, 5); r != 0.5 {
		t.Errorf("five yellows should be worth 0.5, got %v", r)
	}
}

func TestGetActionsSkipsPlayedGuesses(t *testing.T) {
	agreement := NewAgreement([]string{"cigar", "rebut", "soare"})
	universe := game.InitUniverse([]string{"cigar", "rebut"})
	gp := game.InitGamePlay(universe.PossibleWorlds.GetHash())
	if n := len(agreement.GetActions(gp, universe.PossibleWorlds, universe.Actions)); n != 3 {
		t.Fatalf("expected all 3 guesses before any turn, got %d", n)
	}
	gp.AddTurn(Act(gp, game.NewWorld("cigar"), game.NewAction("soare")))
	actions := agreement.GetActions(gp, universe.PossibleWorlds, universe.Actions)
	if len(actions) != 2 || actions[0].String() != "cigar" || actions[1].String() != "rebut" {
		t.Errorf("expected cigar, rebut after soare was played, got %v", actions)
	}
}

// loadPublishedLists reads the original Wordle lists, which are not in this repository: set
// WORDLE_ANSWERS to the 2,315-word answer list and WORDLE_GUESSES to the 12,972-word allowed-guess list.
func loadPublishedLists(t *testing.T) ([]string, []string) {
	answers, guesses := os.Getenv("WORDLE_ANSWERS"), os.Getenv("WORDLE_GUESSES")
	if answers == "" || guesses == "" {
		t.Skip("set WORDLE_ANSWERS and WORDLE_GUESSES to the original Wordle lists to run this test")
	}
	return utils.ReadLinesFromFilePath(answers), utils.ReadLinesFromFilePath(guesses)
}

func classes(answers []string, guess string) map[int]int {
	counts := map[int]int{}
	for _, answer := range answers {
		code, _, _ := Score(answer, guess)
		counts[code]++
	}
	return counts
}

// Published first-guess facts for the original lists (2,315 answers, 12,972 allowed guesses):
// TRACE has 150 response patterns, the most of any guess, and SLATE 147 (real-statistics.com, "Best First
// Two Guesses in Wordle"; caseyho/wordle-solver-benchmarks lists TRACE as the max-partitions opener);
// ARISE and RAISE are the best minimax openers, leaving at most 168 answers (graphallthethings.com/posts/wordle);
// SOARE has the most entropy, 5.89 bits (HKUST Science Focus, "Demystifying Wordle"; caseyho lists SOARE as the max-entropy opener).
func TestPublishedFirstGuessFacts(t *testing.T) {
	answers, allowed := loadPublishedLists(t)
	if len(answers) != 2315 || len(allowed) != 12972 {
		t.Fatalf("expected 2,315 answers and 12,972 allowed guesses, got %d and %d", len(answers), len(allowed))
	}
	if n := len(classes(answers, "trace")); n != 150 {
		t.Errorf("trace: %d classes, want 150", n)
	}
	if n := len(classes(answers, "slate")); n != 147 {
		t.Errorf("slate: %d classes, want 147", n)
	}
	bestParts, bestMax, bestEntropy := 0, math.MaxInt32, 0.0
	partsWord, entropyWord := "", ""
	maxOf := map[string]int{}
	for _, guess := range allowed {
		counts := classes(answers, guess)
		values := []float64{}
		largest := 0
		for _, c := range counts {
			values = append(values, float64(c))
			if c > largest {
				largest = c
			}
		}
		maxOf[guess] = largest
		entropy := mth.CalcEntropy(float64(len(answers)), values)
		if len(counts) > bestParts {
			bestParts, partsWord = len(counts), guess
		}
		if largest < bestMax {
			bestMax = largest
		}
		if entropy > bestEntropy {
			bestEntropy, entropyWord = entropy, guess
		}
	}
	if bestParts != 150 || partsWord != "trace" {
		t.Errorf("most classes: %s with %d, want trace with 150", partsWord, bestParts)
	}
	if bestMax != 168 || maxOf["arise"] != 168 || maxOf["raise"] != 168 {
		t.Errorf("minimax: best largest class %d, arise %d, raise %d; want 168 for both", bestMax, maxOf["arise"], maxOf["raise"])
	}
	if entropyWord != "soare" || bestEntropy < 5.885 || bestEntropy >= 5.895 {
		t.Errorf("max entropy: %s with %.4f bits, want soare with 5.89", entropyWord, bestEntropy)
	}
}

func TestNewAgreementRejectsOtherWordLengths(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected a panic for a nine-letter guess")
		}
	}()
	NewAgreement([]string{"crane", "aardvarks"})
}

// The engine plays Wordle to the end: every game's last action is its answer.
func TestPlayGamesEndsOnTheAnswer(t *testing.T) {
	words := utils.ReadLinesFromFilePath("../../fixtures/five-letter-words.txt")
	sample := words[:40]
	agreement := NewAgreement(words)
	game.ThePredictionsFreeCache = freecache.NewCache(1000 * 1024 * 1024)
	games := play.PlayGames(agreement, game.Motivation{Entropy: 1}, sample, len(sample))
	for i, gameplay := range games {
		last := gameplay.Reactions[len(gameplay.Reactions)-1]
		if last.Action.String() != sample[i] || last.Response.String() != "GGGGG" || last.Reward != 1.0 {
			t.Errorf("game %s ended on %s/%s", sample[i], last.Action.String(), last.Response.String())
		}
	}
}
