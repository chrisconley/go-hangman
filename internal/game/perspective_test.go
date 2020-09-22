package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

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
