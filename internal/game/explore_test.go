package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestBuildPredictions(t *testing.T) {
	lines := []string{"abc", "aab", "bbb"}
	BUniverse = InitUniverse(lines)
	pw := BUniverse.PossibleWorlds
	actions := []Action{"a", "b", "c"}

	actual := ExploreReactions(GamePlay{}, pw, actions, BuildTestAct(1.0))

	test.AssertAlmostEqual(actual.GetEntropy("a"), 0.9183, 4, t)
	test.AssertAlmostEqual(actual.GetEntropy("b"), 0.0000, 4, t)
	test.AssertAlmostEqual(actual.GetEntropy("c"), 0.9183, 4, t)

	test.AssertAlmostEqual(actual.GetReward("a"), 0.6667, 4, t)
	test.AssertAlmostEqual(actual.GetReward("b"), 1.0000, 4, t)
	test.AssertAlmostEqual(actual.GetReward("c"), 0.3333, 4, t)

	test.AssertAlmostEqual(actual.GetInversedMax("a"), 1.0000, 4, t)
	test.AssertAlmostEqual(actual.GetInversedMax("b"), 0.0000, 4, t)
	test.AssertAlmostEqual(actual.GetInversedMax("c"), 1.0000, 4, t)

	test.AssertAlmostEqual(actual.GetPartsCount("a"), 2.0000, 4, t)
	test.AssertAlmostEqual(actual.GetPartsCount("b"), 1.0000, 4, t)
	test.AssertAlmostEqual(actual.GetPartsCount("c"), 2.0000, 4, t)
}

func TestExploreReactionsForAction(t *testing.T) {
	lines := []string{"abc", "aab", "cab"}
	universeLength := len(lines)
	BUniverse = InitUniverse(lines)
	pw := BUniverse.PossibleWorlds

	_, counts, rewards, responseWorlds := ExploreReactionsForAction(
		GamePlay{},
		pw.Indices(),
		Action("c"),
		BuildTestAct(0.222),
	)

	test.AssertAlmostEqual(counts.Get("!"), 1.0, 3, t)
	test.AssertAlmostEqual(counts.Get("c"), 2.0, 3, t)

	test.AssertAlmostEqual(rewards.Get("!"), 0.0, 3, t)
	test.AssertAlmostEqual(rewards.Get("c"), 0.444, 3, t)

	expected := NewPossibleWorlds(universeLength)
	expected.Add(1)
	test.AssertDeepEqual(responseWorlds["!"], expected, t)

	expected = NewPossibleWorlds(universeLength)
	expected.Add(0)
	expected.Add(2)
	test.AssertDeepEqual(responseWorlds["c"], expected, t)
}
