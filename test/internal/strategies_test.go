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

func TestMaximizeEntropy(t *testing.T) {
	// potentialOutcomes := map[string]map[string]int{}
	// potentialOutcomes["a"]["-a-"] = 2
	potentialOutcomes := game.PotentialOutcomes{}
	potentialOutcomes["a"] = game.PossibleResponses{}
	potentialOutcomes["a"]["-a-"] = 4
	potentialOutcomes["a"]["--a"] = 1
	potentialOutcomes["a"]["a--"] = 1
	potentialOutcomes["a"]["!"] = 2
	actual := game.MaximizeEntropy(potentialOutcomes)
	assertAlmostEqual(actual["a"], 1.75, 4, t)
}
