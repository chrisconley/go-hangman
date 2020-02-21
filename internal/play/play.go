package play

import (
	"fmt"
	"math"
	// "regexp"
	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"
)

// ###########################
// ### THE DECISION PROCESS ###
// ###########################
func Explore(perspective game.Perspective, agreement game.Agreement) (predictions game.Predictions) {
	reactions := game.ExploreReactions(perspective.PossibleWorlds, perspective.PossibleActions, agreement)
	predictions = game.BuildPredictions(reactions)
	return predictions
}

func Decide(predictions game.Predictions, motivation game.Motivation) (action game.Action) {
	maxWeight := 0.0
	maxAction := game.Action("")
	for _, prediction := range predictions {
		// fmt.Printf("\nPREDICTION: %+v", prediction)
		weight := math.Pow(prediction.Entropy, motivation.Entropy) * math.Pow(prediction.Reward, motivation.Reward)
		// fmt.Printf("\nWEIGHT: %+v", weight)
		if weight > maxWeight {
			maxWeight = weight
			maxAction = prediction.Action
		}
	}
	// fmt.Printf("\nMAXACTION: %+v", maxAction)
	return maxAction
}

func Act(action game.Action, agreement game.Agreement, truth game.World) (reaction game.Reaction) {
	return agreement(truth, action)
}

func Observe(action game.Action, reaction game.Reaction, perspective *game.Perspective) {
	perspective.UpdateKnowledge(reaction)
	perspective.PossibleActions = hangman.UpdatePossibleActions(reaction, perspective.PossibleActions)

	perspective.PossibleWorlds = hangman.UpdatePossibleWorlds(reaction, perspective.PossibleWorlds)
	// fmt.Printf("\nPOSSIBLE ACTIONS: %+v\n", perspective.PossibleActions)
	// fmt.Printf("\nPOSSIBLE WORLDS: %+v\n", perspective.PossibleWorlds)

}

// ###########################
// ### THE GAME ###
// ###########################

func PlayGame(truth game.World, perspective *game.Perspective) game.World {
	for true {
		playTurn(truth, perspective)
		if len(perspective.PossibleWorlds) == 1 {
			break
		}
	}
	// fmt.Printf("\nGAME OVER: %+v\n", perspective.PossibleWorlds[0])
	return perspective.PossibleWorlds[0].World
}

func playTurn(truth game.World, perspective *game.Perspective) (reaction game.Reaction) {
	predictions := Explore(*perspective, hangman.HangmanAgreement)
	action := Decide(predictions, perspective.Motivation)
	// fmt.Printf("\nACTION: %+v", action)
	reaction = Act(action, hangman.HangmanAgreement, truth)
	// fmt.Printf("\nREACTION: %+v", reaction)
	Observe(action, reaction, perspective)
	return reaction
}

func main() {
	fmt.Printf("\n%+v", game.World("hi"))
}
