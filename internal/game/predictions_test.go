package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestBuildPredictions(t *testing.T) {
	lines := []string{"abc", "aab", "bbb"}
	pw := InitUniverse(lines)
	actions := []Action{"a", "b", "c"}

	ExploreReactions(pw, actions, BuildTestAgreement(1.0))
	actual := BBuildPredictions(actions, pw)

	test.AssertAlmostEqual(actual.GetEntropy("a"), 0.9183, 4, t)
	test.AssertAlmostEqual(actual.GetEntropy("b"), 0.0000, 4, t)
	test.AssertAlmostEqual(actual.GetEntropy("c"), 0.9183, 4, t)

	test.AssertAlmostEqual(actual.GetReward("a"), 0.6667, 4, t)
	test.AssertAlmostEqual(actual.GetReward("b"), 1.0000, 4, t)
	test.AssertAlmostEqual(actual.GetReward("c"), 0.3333, 4, t)
}
