package test

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
)

func assertDeepEqual(x, y interface{}, t *testing.T) {
	if !reflect.DeepEqual(x, y) {
		t.Errorf("\ngot %s\nexp %s", x, y)
	}
}

func assertAlmostEqual(a float64, b float64, places int, t *testing.T) {
	if !(math.Abs(a-b) <= math.Pow10(-places)) {
		t.Errorf("\ngot %f\nexp %f", a, b)
	}
}

func TestLogProbability(t *testing.T) {
	assertAlmostEqual(game.LogProbability(0.0), 0.0, 4, t)
	assertAlmostEqual(game.LogProbability(0.5), 0.5, 4, t)
	assertAlmostEqual(game.LogProbability(1.0), 0.0, 4, t)
	assertAlmostEqual(game.LogProbability(0.01), 0.0664, 4, t)
	assertAlmostEqual(game.LogProbability(0.99), 0.0144, 4, t)
}

func TestEntropy(t *testing.T) {
	reactionCountsByAction := game.ReactionCountsByAction{}
	reactionCountsByAction["a"] = game.ReactionCounts{ResponseCounts: game.ResponseCounts{}, Reward: game.Reward{0, 0}}
	reactionCountsByAction["a"].ResponseCounts["-a-"] = &game.Counts{Weight: 0, Count: 4}
	reactionCountsByAction["a"].ResponseCounts["--a"] = &game.Counts{Weight: 0, Count: 1}
	reactionCountsByAction["a"].ResponseCounts["a--"] = &game.Counts{Weight: 0, Count: 1}
	reactionCountsByAction["a"].ResponseCounts["!"] = &game.Counts{Weight: 0, Count: 2}
	actual := game.BuildPredictions(reactionCountsByAction)
	assertAlmostEqual(actual[game.Action("a")].Entropy, 1.75, 4, t)
}

func testAgreement(w game.World, a game.Action) (r game.Reaction) {
	if strings.Contains(string(w), string(a)) {
		return game.Reaction{Response: game.Response(a), Reward: 1.0}
	} else {
		return game.Reaction{Response: "!", Reward: 0.0}
	}
}

func TestExploreReactions(t *testing.T) {
	worlds := game.PossibleWorlds{game.PossibleWorld{"abc", 0.5}, game.PossibleWorld{"aab", 0.5}}
	actions := game.PossibleActions{"a", "b", "c"}
	actual := game.ExploreReactions(worlds, actions, testAgreement)
	assertDeepEqual(actual["a"].ResponseCounts["a"].Count, 2, t)
	assertDeepEqual(actual["c"].ResponseCounts["c"].Count, 1, t)
	assertDeepEqual(actual["c"].ResponseCounts["!"].Count, 1, t)
}
