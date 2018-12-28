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
	reactionCountsByAction.SetResponseCounts("a", "-a-", 0, 4)
	reactionCountsByAction.SetResponseCounts("a", "--a", 0, 1)
	reactionCountsByAction.SetResponseCounts("a", "a--", 0, 1)
	reactionCountsByAction.SetResponseCounts("a", "!", 0, 2)
	actual := game.BuildPredictions(reactionCountsByAction)
	assertAlmostEqual(actual.GetEntropy("a"), 1.75, 4, t)
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
	assertDeepEqual(actual.GetResponseCount("a", "a"), 2, t)
	assertDeepEqual(actual.GetResponseCount("c", "c"), 1, t)
	assertDeepEqual(actual.GetResponseCount("c", "!"), 1, t)
}

func TestIntegration(t *testing.T) {
	worlds := game.PossibleWorlds{
		game.PossibleWorld{"abc", 0.5},
		game.PossibleWorld{"aab", 0.5},
		game.PossibleWorld{"bbb", 0.5},
	}
	actions := game.PossibleActions{"a", "b", "c"}

	reactions := game.ExploreReactions(worlds, actions, testAgreement)
	actual := game.BuildPredictions(reactions)

	assertAlmostEqual(actual.GetEntropy("a"), 0.9183, 4, t)
	assertAlmostEqual(actual.GetEntropy("b"), 0.0000, 4, t)
	assertAlmostEqual(actual.GetEntropy("c"), 0.9183, 4, t)
}
