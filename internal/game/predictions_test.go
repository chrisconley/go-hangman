package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestBuildPredictions(t *testing.T) {
	worlds := Worlds{
		World("abc"),
		World("aab"),
		World("bbb"),
	}
	aa := BuildAvailableActions(worlds)

	reactions := ExploreReactions(worlds, aa.AllActions(), TestAgreement)
	actual := BuildPredictions(reactions)

	test.AssertAlmostEqual(actual.GetEntropy("a"), 0.9183, 4, t)
	test.AssertAlmostEqual(actual.GetEntropy("b"), 0.0000, 4, t)
	test.AssertAlmostEqual(actual.GetEntropy("c"), 0.9183, 4, t)
}
