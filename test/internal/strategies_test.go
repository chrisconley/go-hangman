package test

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/test"
)

func TestLogProbability(t *testing.T) {
	test.AssertAlmostEqual(game.LogProbability(0.0), 0.0, 4, t)
	test.AssertAlmostEqual(game.LogProbability(0.5), 0.5, 4, t)
	test.AssertAlmostEqual(game.LogProbability(1.0), 0.0, 4, t)
	test.AssertAlmostEqual(game.LogProbability(0.01), 0.0664, 4, t)
	test.AssertAlmostEqual(game.LogProbability(0.99), 0.0144, 4, t)
}

func TestEntropy(t *testing.T) {
	reactionCountsByAction := game.ReactionCountsByAction{}
	reactionCountsByAction.SetResponseCounts("a", "-a-", 0, 4)
	reactionCountsByAction.SetResponseCounts("a", "--a", 0, 1)
	reactionCountsByAction.SetResponseCounts("a", "a--", 0, 1)
	reactionCountsByAction.SetResponseCounts("a", "!", 0, 2)
	actual := game.BuildPredictions(reactionCountsByAction)
	test.AssertAlmostEqual(actual.GetEntropy("a"), 1.75, 4, t)
}

func TestExploreReactions(t *testing.T) {
	worlds := game.PossibleWorlds{game.PossibleWorld{"abc", 0.5}, game.PossibleWorld{"aab", 0.5}}
	actions := game.PossibleActions{"a", "b", "c"}
	actual := game.ExploreReactions(worlds, actions, test.TestAgreement)
	test.AssertDeepEqual(actual.GetResponseCount("a", "a"), 2, t)
	test.AssertDeepEqual(actual.GetResponseCount("c", "c"), 1, t)
	test.AssertDeepEqual(actual.GetResponseCount("c", "!"), 1, t)
}

func TestIntegration(t *testing.T) {
	worlds := game.PossibleWorlds{
		game.PossibleWorld{"abc", 0.5},
		game.PossibleWorld{"aab", 0.5},
		game.PossibleWorld{"bbb", 0.5},
	}
	actions := game.PossibleActions{"a", "b", "c"}

	reactions := game.ExploreReactions(worlds, actions, test.TestAgreement)
	actual := game.BuildPredictions(reactions)

	test.AssertAlmostEqual(actual.GetEntropy("a"), 0.9183, 4, t)
	test.AssertAlmostEqual(actual.GetEntropy("b"), 0.0000, 4, t)
	test.AssertAlmostEqual(actual.GetEntropy("c"), 0.9183, 4, t)
}
