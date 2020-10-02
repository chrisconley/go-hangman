package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestDecide(t *testing.T) {
	predictions := Predictions{
		"a": Prediction{Action: "a", Entropy: 1.75, Reward: 0.9},
		"b": Prediction{Action: "b", Entropy: 2.50, Reward: 0.1},
	}

	// All Entropy
	motivation := Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
	actual := Decide(predictions, motivation)
	test.AssertDeepEqual(actual, predictions["b"], t)

	// Weighted Reward with Entropy
	motivation = Motivation{Reward: 0.75, Entropy: 0.25, Time: 0.0}
	actual = Decide(predictions, motivation)
	test.AssertDeepEqual(actual, predictions["a"], t)
}

func TestDecideDeterminism(t *testing.T) {
	predictions := Predictions{
		"a": Prediction{Action: "a", Reward: 1, Entropy: 0},
		"b": Prediction{Action: "b", Reward: 1, Entropy: 0},
		"c": Prediction{Action: "c", Reward: 0, Entropy: 0},
		"d": Prediction{Action: "d", Reward: 0.5, Entropy: 1},
		"e": Prediction{Action: "e", Reward: 1, Entropy: 0},
		"f": Prediction{Action: "f", Reward: 0, Entropy: 0},
		"g": Prediction{Action: "g", Reward: 0, Entropy: 0},
		"h": Prediction{Action: "h", Reward: 0, Entropy: 0},
		"i": Prediction{Action: "i", Reward: 1, Entropy: 0},
		"j": Prediction{Action: "j", Reward: 0, Entropy: 0},
		"k": Prediction{Action: "k", Reward: 0, Entropy: 0},
		"l": Prediction{Action: "l", Reward: 0, Entropy: 0},
		"m": Prediction{Action: "m", Reward: 0, Entropy: 0},
		"n": Prediction{Action: "n", Reward: 0, Entropy: 0},
		"o": Prediction{Action: "o", Reward: 1, Entropy: 0},
		"p": Prediction{Action: "p", Reward: 0, Entropy: 0},
		"q": Prediction{Action: "q", Reward: 0, Entropy: 0},
		"r": Prediction{Action: "r", Reward: 0, Entropy: 0},
		"s": Prediction{Action: "s", Reward: 0.5, Entropy: 1},
		"t": Prediction{Action: "t", Reward: 0.5, Entropy: 1},
		"u": Prediction{Action: "u", Reward: 1, Entropy: 0},
		"v": Prediction{Action: "v", Reward: 0, Entropy: 0},
		"w": Prediction{Action: "w", Reward: 0, Entropy: 0},
		"x": Prediction{Action: "x", Reward: 0.5, Entropy: 1},
		"y": Prediction{Action: "y", Reward: 0, Entropy: 0},
		"z": Prediction{Action: "z", Reward: 0, Entropy: 0},
	}

	motivation := Motivation{Reward: 0.5, Entropy: 0.25, Time: 0.25}
	for i := 0; i < 100; i++ {
		actual := Decide(predictions, motivation)
		test.AssertDeepEqual(actual, predictions["d"], t)
	}
}

func TestDecideLexicallySortsOnTie(t *testing.T) {
	predictions := Predictions{
		"b": Prediction{Action: "b", Entropy: 1.0, ActionMatchesAPossibleWorld: true},
		"a": Prediction{Action: "a", Entropy: 1.0, ActionMatchesAPossibleWorld: true},
	}

	// All Entropy
	motivation := Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
	actual := Decide(predictions, motivation)
	test.AssertDeepEqual(actual, predictions["a"], t)

	predictions = Predictions{
		"b": Prediction{Action: "b", Entropy: 1.0, ActionMatchesAPossibleWorld: false},
		"a": Prediction{Action: "a", Entropy: 1.0, ActionMatchesAPossibleWorld: false},
	}

	// All Entropy
	motivation = Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
	actual = Decide(predictions, motivation)
	test.AssertDeepEqual(actual, predictions["a"], t)
}
func TestDecideValidSortsOnTie(t *testing.T) {
	predictions := Predictions{
		"b": Prediction{Action: "b", Entropy: 1.0, ActionMatchesAPossibleWorld: false},
		"c": Prediction{Action: "c", Entropy: 1.0, ActionMatchesAPossibleWorld: true},
		"a": Prediction{Action: "a", Entropy: 1.0, ActionMatchesAPossibleWorld: false},
	}

	// All Entropy
	motivation := Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
	actual := Decide(predictions, motivation)
	test.AssertDeepEqual(actual, predictions["c"], t)

	predictions = Predictions{
		"b": Prediction{Action: "b", Entropy: 1.0, ActionMatchesAPossibleWorld: true},
		"c": Prediction{Action: "c", Entropy: 1.0, ActionMatchesAPossibleWorld: true},
		"a": Prediction{Action: "a", Entropy: 1.0, ActionMatchesAPossibleWorld: false},
	}

	// All Entropy
	motivation = Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
	actual = Decide(predictions, motivation)
	test.AssertDeepEqual(actual, predictions["b"], t)
}
