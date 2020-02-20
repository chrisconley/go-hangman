package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestBuildPredictions(t *testing.T) {
	reactionCounter := ReactionCounter{}
	reactionCounter.Set(Reaction{Action: "a", Response: "-a-"}, 4.0, 4)
	reactionCounter.Set(Reaction{Action: "a", Response: "--a"}, 1.0, 1)
	reactionCounter.Set(Reaction{Action: "a", Response: "a--"}, 1.0, 1)
	reactionCounter.Set(Reaction{Action: "a", Response: "!"}, 0.0, 2)
	actual := BuildPredictions(reactionCounter)
	test.AssertAlmostEqual(actual.GetEntropy("a"), 1.75, 4, t)
	test.AssertAlmostEqual(actual.GetReward("a"), 0.75, 4, t)
}
