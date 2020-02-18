package test

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
)

func TestAgreement(w game.World, a game.Action) (r game.Reaction) {
	if strings.Contains(string(w), string(a)) {
		return game.Reaction{Response: game.Response(a), Reward: 1.0}
	} else {
		return game.Reaction{Response: "!", Reward: 0.0}
	}
}

func AssertDeepEqual(x, y interface{}, t *testing.T) {
	if !reflect.DeepEqual(x, y) {
		t.Errorf("\ngot %s\nexp %s", x, y)
	}
}

func AssertAlmostEqual(a float64, b float64, places int, t *testing.T) {
	if !(math.Abs(a-b) <= math.Pow10(-places)) {
		t.Errorf("\ngot %f\nexp %f", a, b)
	}
}
