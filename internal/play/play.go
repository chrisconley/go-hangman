package play

import (
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
	if len(predictions) == 1 {
		panic("There is only one prediction for the decision.")
	}
	for _, prediction := range predictions {
		weight := math.Pow(prediction.Reward, motivation.Reward) *
			math.Pow(prediction.Entropy, motivation.Entropy) *
			math.Pow(prediction.InversedMax, motivation.MiniMax)

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

func Act(action game.Action, act game.Act, truth game.World) (reaction game.Reaction) {
	return act(truth, action)
}

func Observe(action game.Action, reaction game.Reaction, possibleWorlds game.BPossibleWorlds) game.BPossibleWorlds {
	hash := game.BuildActionWorldHash(action, possibleWorlds.GetHash())
	response := string(reaction.Response)
	return game.BBPossibleWorldsStore[hash][response]
}

// ###########################
// ### THE GAME ###
// ###########################
func PlayGames(agreement game.Agreement, motivation game.Motivation, lines []string, limit int) []int {
	game.InitUniverse(lines)
	// defer profile.Start().Stop()
	turns := []int{}
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

func PlayGame(truth game.World, perspective game.Perspective) (int, game.World) {
	possibleWorlds := game.BUniverse.BPossibleWorlds
	turns := 0
	for true {
		turns += 1
		_, possibleWorlds = playTurn(truth, perspective, possibleWorlds)
		remainingWorldIndices := possibleWorlds.Indices()
		if len(remainingWorldIndices) == 1 {
			break
		}
	}
	remainingWorldIndices := possibleWorlds.Indices()
	lastWorldIndex := remainingWorldIndices[0]
	return turns, game.World(game.BUniverse.Worlds[lastWorldIndex].Name)
}

func playTurn(truth game.World, perspective game.Perspective, possibleWorlds game.BPossibleWorlds) (game.Reaction, game.BPossibleWorlds) {
	predictions := Explore(perspective.Agreement, possibleWorlds)
	action := Decide(predictions, perspective.Motivation)
	reaction := Act(action, perspective.Agreement.Act, truth)
	possibleWorlds = Observe(action, reaction, possibleWorlds)
	return reaction, possibleWorlds
}
