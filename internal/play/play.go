package play

import (
	"math"

	"github.com/chrisconley/go-hangman/internal/game"
)

// ###########################
// ### THE DECISION PROCESS ###
// ###########################
func Explore(perspective *game.Perspective, agreement game.Agreement) (predictions game.Predictions) {
	actions := perspective.PlausibleActions(perspective.PossibleWorlds)
	game.ExploreReactions(perspective.BPossibleWorlds, actions, agreement)
	predictions = game.BBuildPredictions(actions, perspective.BPossibleWorlds)
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

func Observe(action game.Action, reaction game.Reaction, perspective *game.Perspective) {
	perspective.UpdateKnowledge(reaction)
	hash := game.BuildActionWorldHash(action, perspective.BPossibleWorlds)
	response := string(reaction.Response)
	perspective.BPossibleWorlds = game.BBPossibleWorldsStore[hash][response]
}

// ###########################
// ### THE GAME ###
// ###########################

func PlayGame(truth game.World, perspective *game.Perspective, memory game.DecisionMemory) game.World {
	gameState := game.GameState{}
	for true {
		reaction, possibleWorlds := playTurn(truth, perspective)
		gameState.Add(reaction)
		memory.Add(gameState, possibleWorlds)
		remainingWorldIndices := perspective.BPossibleWorlds.ToNums()
		if len(remainingWorldIndices) == 1 {
			break
		}
	}
	remainingWorldIndices := perspective.BPossibleWorlds.ToNums()
	lastWorldIndex := remainingWorldIndices[0]
	return game.World(game.BUniverse[lastWorldIndex].Name)
}

func playTurn(truth game.World, perspective *game.Perspective) (game.Reaction, game.PossibleWorlds) {
	predictions := Explore(perspective, perspective.Agreement)
	action := Decide(predictions, perspective.Motivation)
	reaction := Act(action, perspective.Agreement, truth)
	Observe(action, reaction, perspective)
	return reaction, perspective.PossibleWorlds
}
