package play

import (
	"math"

	"github.com/chrisconley/go-hangman/internal/game"
)

// ###########################
// ### THE DECISION PROCESS ###
// ###########################
func Explore(perspective *game.Perspective, agreement game.Agreement, possibleWorlds game.BPossibleWorlds) (predictions game.Predictions) {
	actions := perspective.PlausibleActions(perspective.PossibleWorlds)
	game.ExploreReactions(possibleWorlds, actions, agreement)
	predictions = game.BBuildPredictions(actions, possibleWorlds)
	return predictions
}

func Decide(predictions game.Predictions, motivation game.Motivation) game.Action {
	maxWeight := 0.0
	maxAction := game.Action("")
	if len(predictions) == 1 {
		panic("There is only one prediction for the decision.")
	}
	for _, prediction := range predictions {
		weight := math.Pow(prediction.Entropy, motivation.Entropy) * math.Pow(prediction.Reward, motivation.Reward)

		// Decide whether or not to swap out maxAction with the new action
		if weight > maxWeight {
			// If there is a clear winner, use it
			maxWeight = weight
			maxAction = prediction.Action
		} else if weight == maxWeight && prediction.Action < maxAction {
			// If the new weight is equal to current max
			// AND the new action is lexicographically before maxAction, use it
			maxAction = prediction.Action
		}
	}
	return maxAction
}

func Act(action game.Action, agreement game.Agreement, truth game.World) (reaction game.Reaction) {
	return agreement(truth, action)
}

func Observe(action game.Action, reaction game.Reaction, perspective *game.Perspective, possibleWorlds game.BPossibleWorlds) game.BPossibleWorlds {
	perspective.UpdateKnowledge(reaction)
	hash := game.BuildActionWorldHash(action, possibleWorlds)
	response := string(reaction.Response)
	return game.BBPossibleWorldsStore[hash][response]
}

// ###########################
// ### THE GAME ###
// ###########################

func PlayGame(truth game.World, perspective *game.Perspective) game.World {
	possibleWorlds := perspective.BPossibleWorlds
	for true {
		_, possibleWorlds = playTurn(truth, perspective, possibleWorlds)
		remainingWorldIndices := possibleWorlds.ToNums()
		if len(remainingWorldIndices) == 1 {
			break
		}
	}
	remainingWorldIndices := possibleWorlds.ToNums()
	lastWorldIndex := remainingWorldIndices[0]
	return game.World(game.BUniverse[lastWorldIndex].Name)
}

func playTurn(truth game.World, perspective *game.Perspective, possibleWorlds game.BPossibleWorlds) (game.Reaction, game.BPossibleWorlds) {
	predictions := Explore(perspective, perspective.Agreement, possibleWorlds)
	action := Decide(predictions, perspective.Motivation)
	reaction := Act(action, perspective.Agreement, truth)
	possibleWorlds = Observe(action, reaction, perspective, possibleWorlds)
	return reaction, possibleWorlds
}
