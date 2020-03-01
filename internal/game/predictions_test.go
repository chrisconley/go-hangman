package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestBuildPredictions(t *testing.T) {
	worlds := PossibleWorlds{
		PossibleWorld{"abc", 0.5},
		PossibleWorld{"aab", 0.5},
		PossibleWorld{"bbb", 0.5},
	}
	actions := PossibleActions{"a", "b", "c"}

	reactions := ExploreReactions(worlds, actions, TestAgreement)
	actual := BuildPredictions(reactions)

	test.AssertAlmostEqual(actual.GetEntropy("a"), 0.9183, 4, t)
	test.AssertAlmostEqual(actual.GetEntropy("b"), 0.0000, 4, t)
	test.AssertAlmostEqual(actual.GetEntropy("c"), 0.9183, 4, t)
}
