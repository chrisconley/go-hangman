package main

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"
	"github.com/chrisconley/go-hangman/test"
)

func TestDecide(t *testing.T) {
	predictions := game.Predictions{
		"a": game.Prediction{Action: "a", Entropy: 1.75, Reward: 0.9},
		"b": game.Prediction{Action: "b", Entropy: 2.50, Reward: 0.1},
	}

	// All Entropy
	motivation := game.Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
	actual := Decide(predictions, motivation)
	test.AssertDeepEqual(actual, game.Action("b"), t)

	// Weighted Reward with Entropy
	motivation = game.Motivation{Reward: 0.75, Entropy: 0.25, Time: 0.0}
	actual = Decide(predictions, motivation)
	test.AssertDeepEqual(actual, game.Action("a"), t)
}

func TestAct(t *testing.T) {
	truth := game.World("abc")
	actual := Act(game.Action("a"), test.TestAgreement, truth)
	test.AssertDeepEqual(actual, game.Reaction{Response: "a", Reward: 1.0}, t)
}

func TestObserve(t *testing.T) {
	perspective := &game.Perspective{}
	perspective.Knowledge = game.Knowledge{}

	reaction := game.Reaction{Response: "a", Reward: 1.0}
	Observe(game.Action("a"), reaction, perspective)
	test.AssertDeepEqual(perspective.Knowledge, game.Knowledge{reaction}, t)
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
	perspective.Agreement = hangman.HangmanAgreement
	perspective.Motivation = game.Motivation{0.5, 0.25, 0.25}

	truth := game.World("bbb")
	actual := PlayGame(truth, perspective)
	test.AssertDeepEqual(actual, truth, t)
}
