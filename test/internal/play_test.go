package test

import (
	"strings"
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
)

func TestDecide(t *testing.T) {
	predictions := game.Predictions{
		"a": game.PredictedActionUtility{Action: "a", Entropy: 1.75, Reward: 0.9},
		"b": game.PredictedActionUtility{Action: "b", Entropy: 2.50, Reward: 0.1},
	}

	// All Entropy
	motivation := game.Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
	actual := game.Decide(predictions, motivation)
	assertDeepEqual(actual, game.Action("b"), t)

	// Weighted Reward with Entropy
	motivation = game.Motivation{Reward: 0.75, Entropy: 0.25, Time: 0.0}
	actual = game.Decide(predictions, motivation)
	assertDeepEqual(actual, game.Action("a"), t)
}

func testAgreement(w game.World, a game.Action) (r game.Reaction) {
	if strings.Contains(string(w), string(a)) {
		return game.Reaction{Response: game.Response(a), Reward: 1.0}
	} else {
		return game.Reaction{Response: "!", Reward: 0.0}
	}
}

func TestAct(t *testing.T) {
	truth := game.World("abc")
	actual := game.Act(game.Action("a"), testAgreement, truth)
	assertDeepEqual(actual, game.Reaction{Response: "a", Reward: 1.0}, t)
}

func TestObserve(t *testing.T) {
	perspective := &game.Perspective{}
	perspective.Knowledge = game.Knowledge{}

	reaction := game.Reaction{Response: "a", Reward: 1.0}
	game.Observe(game.Action("a"), reaction, perspective)
	assertDeepEqual(perspective.Knowledge, game.Knowledge{reaction}, t)
}

func TestPlayGame(t *testing.T) {
	perspective := &game.Perspective{}
	perspective.Knowledge = game.Knowledge{}
	perspective.PossibleWorlds = game.PossibleWorlds{
		game.PossibleWorld{"abc", 0.5},
		game.PossibleWorld{"aab", 0.5},
		game.PossibleWorld{"bbb", 0.5},
	}
	perspective.PossibleActions = game.PossibleActions{"a", "b", "c"}
	perspective.Agreement = game.HangmanAgreement
	perspective.Motivation = game.Motivation{0.5, 0.25, 0.25}

	truth := game.World("bbb")
	actual := game.PlayGame(truth, perspective)
	assertDeepEqual(actual, truth, t)
}