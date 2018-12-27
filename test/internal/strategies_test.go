package test

import (
	"math"
	"reflect"
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
	reactionCountsByAction["a"] = game.ReactionCounts{}
	reactionCountsByAction["a"]["-a-"] = 4
	reactionCountsByAction["a"]["--a"] = 1
	reactionCountsByAction["a"]["a--"] = 1
	reactionCountsByAction["a"]["!"] = 2
	actual := game.BuildPredictions(reactionCountsByAction)
	assertAlmostEqual(actual[game.Action("a")].Entropy, 1.75, 4, t)
}
