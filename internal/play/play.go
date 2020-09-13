package play

import (
	"fmt"
	"math"

	// "regexp"
	"github.com/chrisconley/go-hangman/internal/game"
)

// ###########################
// ### THE DECISION PROCESS ###
// ###########################
func Explore(perspective *game.Perspective, agreement game.Agreement) (predictions game.Predictions) {
	reactions := game.ExploreReactions(perspective.PossibleWorlds, perspective.PlausibleActions(perspective.PossibleWorlds), agreement)
	predictions = game.BuildPredictions(reactions)
	perspective.UpdateReactionExplorer(reactions)
	// fmt.Printf("\nEXPLORE: REACTION EXPLORER: %+v\n", perspective.ReactionExplorer)
	return predictions
}

func Decide(predictions game.Predictions, motivation game.Motivation) (action game.Action) {
	maxWeight := 0.0
	maxAction := game.Action("")
	if len(predictions) == 1 {
		panic("There is only one prediction for the decision.")
	}
	for _, prediction := range predictions {
		// fmt.Printf("\nPREDICTION: %+v", prediction)
		weight := math.Pow(prediction.Entropy, motivation.Entropy) * math.Pow(prediction.Reward, motivation.Reward)

		// Decide whether or not to swap out maxAction with the new action
		if weight > maxWeight {
			// If there is a clear winner, use it
			maxWeight = weight
			maxAction = prediction.Action
		} else if weight == maxWeight && action < prediction.Action {
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
	// fmt.Printf("\nOBSERVE: REACTION EXPLORER: %+v\n", perspective.ReactionExplorer)
	// fmt.Printf("\nOBSERVE: REACTION: %+v\n", reaction)
	perspective.UpdateKnowledge(reaction)
	perspective.PossibleWorlds = perspective.ReactionExplorer.Get(reaction.Action, reaction.Response)
	// fmt.Printf("\nOBSERVE: POSSIBLE WORLDS: %+v\n", perspective.PossibleWorlds)
	// fmt.Printf("\nPOSSIBLE WORLDS: %+v\n", perspective.PossibleWorlds)

}

// ###########################
// ### THE GAME ###
// ###########################

func PlayGame(truth game.World, perspective *game.Perspective, memory game.DecisionMemory) game.World {
	// TODO: Validate that truth is within perspective.PossibleWorlds
	gameState := game.GameState{}
	for true {
		reaction, possibleWorlds := playTurn(truth, perspective)
		gameState.Add(reaction)
		memory.Add(gameState, possibleWorlds)
		if perspective.PossibleWorlds.Len() == 1 {
			break
		}
	}
	// fmt.Printf("\nGAME OVER: %+v\n", perspective.PossibleWorlds[0])
	return perspective.PossibleWorlds.Worlds[0]
}

func playTurn(truth game.World, perspective *game.Perspective) (game.Reaction, game.PossibleWorlds) {
	predictions := Explore(perspective, perspective.Agreement)
	action := Decide(predictions, perspective.Motivation)
	// fmt.Printf("\nACTION: %+v\n", action)
	reaction := Act(action, perspective.Agreement, truth)
	// fmt.Printf("\nREACTION: %+v", reaction)
	Observe(action, reaction, perspective)
	return reaction, perspective.PossibleWorlds
}

func main() {
	fmt.Printf("\n%+v", game.World("hi"))
}
