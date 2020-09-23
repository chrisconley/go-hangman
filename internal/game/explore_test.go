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

func TestExploreReactionsForAction(t *testing.T) {
	lines := []string{"abc", "aab", "cab"}
	universeLength := len(lines)
	pw := InitUniverse(lines)

	counts, rewards, responseWorlds := ExploreReactionsForAction(
		pw,
		Action("c"),
		BuildTestAgreement(0.222),
	)

	test.AssertAlmostEqual(counts.Get("!"), 1.0, 3, t)
	test.AssertAlmostEqual(counts.Get("c"), 2.0, 3, t)

	test.AssertAlmostEqual(rewards.Get("!"), 0.0, 3, t)
	test.AssertAlmostEqual(rewards.Get("c"), 0.444, 3, t)

	expected := NewPossibleWorlds(universeLength)
	expected.SetBit(1)
	test.AssertDeepEqual(responseWorlds["!"], expected, t)

	expected = NewPossibleWorlds(universeLength)
	expected.SetBit(0)
	expected.SetBit(2)
	test.AssertDeepEqual(responseWorlds["c"], expected, t)
}
