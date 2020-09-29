package play

import (
	"fmt"
	"math"

	"github.com/chrisconley/go-hangman/internal/game"
)

// ###########################
// ### THE DECISION PROCESS ###
// ###########################
func Explore(agreement game.Agreement, possibleWorlds game.BPossibleWorlds) (predictions game.Predictions) {
	actions := agreement.GetActions(
		game.BUniverse.PossibleWorlds,
		game.BUniverse.Actions,
	)
	game.ExploreReactions(possibleWorlds, actions, agreement.Act)
	predictions = game.BBuildPredictions(actions, possibleWorlds)
	return predictions
}

func Decide(predictions game.Predictions, motivation game.Motivation) game.Action {
	maxWeight := 0.0
	maxAction := game.Action("")
	maxIsValid := false
	if len(predictions) == 1 {
		panic("There is only one prediction for the decision.")
	}
	for _, prediction := range predictions {
		weight := math.Pow(prediction.Reward, motivation.Reward) *
			math.Pow(prediction.Entropy, motivation.Entropy) *
			math.Pow(prediction.InversedMax, motivation.MiniMax)

		isValid := prediction.ActionMatchesAPossibleWorld

		// Decide whether or not to swap out maxAction with the new action
		if weight == maxWeight && isValid == maxIsValid && prediction.Action < maxAction {
			// 	// If the new weight is equal to current max
			// 	// AND the new action is lexicographically before maxAction, use it
			maxAction = prediction.Action
		} else if weight == maxWeight && isValid && maxIsValid == false {
			maxAction = prediction.Action
			maxIsValid = isValid
		} else if weight > maxWeight {
			// If there is a clear winner, use it
			maxWeight = weight
			maxAction = prediction.Action
			maxIsValid = isValid
		}
	}
	fmt.Printf("MAXACTION: %+v\n", predictions[maxAction])
	return maxAction
}

func Act(action game.Action, act game.Act, truth game.World) game.Reaction {
	reaction := act(truth, action)
	fmt.Printf("REACTION: %+v\n", reaction)
	return reaction
}

func Observe(action game.Action, reaction game.Reaction, possibleWorlds game.BPossibleWorlds) game.BPossibleWorlds {
	hash := game.BuildActionWorldHash(action, possibleWorlds.GetHash())
	response := string(reaction.Response)
	return game.BBPossibleWorldsStore[hash][response]
}

// ###########################
// ### THE GAME ###
// ###########################
func PlayGames(agreement game.Agreement, motivation game.Motivation, lines []string, limit int) []game.Reactions {
	game.InitUniverse(lines)
	// defer profile.Start().Stop()
	turns := []game.Reactions{}
	for j := 0; j < limit; j++ {
		perspective := game.NewPerspective(
			motivation,
			agreement,
		)

		btruth := game.BUniverse.Worlds[j]
		truth := game.World(btruth.Name)
		turnCount, _ := PlayGame(truth, perspective)
		// fmt.Printf("TRUTH: %s | ", truth)
		// fmt.Printf("TURNS: %d\n", turnCount)
		turns = append(turns, turnCount)
	}

	return turns
}

func PlayGame(truth game.World, perspective game.Perspective) (game.Reactions, game.World) {
	possibleWorlds := game.BUniverse.BPossibleWorlds
	reactions := game.Reactions{}
	reaction := game.Reaction{}
	for true {
		// fmt.Printf("POSSIBLE WORLDS: %+v\n", possibleWorlds)
		reaction, possibleWorlds = playTurn(truth, perspective, possibleWorlds)
		reactions.Turns = append(reactions.Turns, reaction)
		reactions.TurnCount++
		remainingWorldIndices := possibleWorlds.Indices()
		if len(remainingWorldIndices) == 1 {
			break
		}
	}
	remainingWorldIndices := possibleWorlds.Indices()
	lastWorldIndex := remainingWorldIndices[0]
	return reactions, game.World(game.BUniverse.Worlds[lastWorldIndex].Name)
}

func playTurn(truth game.World, perspective game.Perspective, possibleWorlds game.BPossibleWorlds) (game.Reaction, game.BPossibleWorlds) {
	predictions := Explore(perspective.Agreement, possibleWorlds)
	action := Decide(predictions, perspective.Motivation)
	reaction := Act(action, perspective.Agreement.Act, truth)
	possibleWorlds = Observe(action, reaction, possibleWorlds)
	return reaction, possibleWorlds
}
