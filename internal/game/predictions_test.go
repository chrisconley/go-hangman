package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestEntropy(t *testing.T) {
	reactionCountsByAction := ReactionCountsByAction{}
	reactionCountsByAction.SetResponseCounts("a", "-a-", 0, 4)
	reactionCountsByAction.SetResponseCounts("a", "--a", 0, 1)
	reactionCountsByAction.SetResponseCounts("a", "a--", 0, 1)
	reactionCountsByAction.SetResponseCounts("a", "!", 0, 2)
	actual := BuildPredictions(reactionCountsByAction)
	test.AssertAlmostEqual(actual.GetEntropy("a"), 1.75, 4, t)
}

func TestExploreReactions(t *testing.T) {
	worlds := PossibleWorlds{PossibleWorld{"abc", 0.5}, PossibleWorld{"aab", 0.5}}
	actions := PossibleActions{"a", "b", "c"}
	actual := ExploreReactions(worlds, actions, TestAgreement)
	test.AssertDeepEqual(actual.GetResponseCount("a", "a"), 2, t)
	test.AssertDeepEqual(actual.GetResponseCount("c", "c"), 1, t)
	test.AssertDeepEqual(actual.GetResponseCount("c", "!"), 1, t)
}

func TestIntegration(t *testing.T) {
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
