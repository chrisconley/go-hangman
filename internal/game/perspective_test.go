package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestPerspectivePlausibleActions(t *testing.T) {
	remainingWorlds := BuildPossibleWorlds([]string{"abc", "aab"})
	perspective := &Perspective{}
	perspective.AvailableActions = BuildAvailableActions(remainingWorlds)
	perspective.ActionFilter = TestActionFilter

	actual := perspective.PlausibleActions(remainingWorlds)
	test.AssertDeepEqual(actual, []Action{"hi"}, t)
}

func TestExploreReactionsForAction(t *testing.T) {
	lines := []string{"abc", "aab", "cab"}
	universeLength := len(lines)
	pw := InitUniverse(lines)

	counts, rewards, responseWorlds := ExploreReactionsForAction(
		pw,
		Action("c"),
		TestAgreement,
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
