package play

import (
	"fmt"
	"os"

	"github.com/chrisconley/go-hangman/internal/game"
)

func logPredictions(predictions game.Predictions) {
	if os.Getenv("DEBUGPREDICTIONS") != "" {
		for _, prediction := range predictions {
			fmt.Printf("PREDICTION: %+v\n", prediction)
		}
	}
}

func logDecision(decision game.Prediction) {
	if os.Getenv("DEBUGPREDICTIONS") != "" {
		fmt.Printf("DECISION: %+v\n", decision)
	}
}

func logReaction(reaction game.Reaction) {
	if os.Getenv("DEBUGPREDICTIONS") != "" {
		fmt.Printf("REACTION: %+v\n", reaction)
	}
}

// ###########################
// ### THE DECISION PROCESS ###
// ###########################
func Explore(gameplay game.GamePlay, agreement game.Agreement, possibleWorlds game.PossibleWorlds) (predictions game.Predictions) {
	actions := agreement.GetActions(
		possibleWorlds,
		game.BUniverse.Actions,
	)
	game.ExploreReactions(gameplay, possibleWorlds, actions, agreement.Act)
	predictions = game.BBuildPredictions(actions, possibleWorlds)
	logPredictions(predictions)
	return predictions
}

func Decide(predictions game.Predictions, motivation game.Motivation) game.Action {
	prediction := game.Decide(predictions, motivation)
	logDecision(prediction)
	return prediction.Action
}

func Act(gameplay game.GamePlay, action game.Action, act game.Act, truth game.World) game.Reaction {
	reaction := act(gameplay, truth, action)
	logReaction(reaction)
	return reaction
}

func Observe(action game.Action, reaction game.Reaction, possibleWorlds game.PossibleWorlds) game.PossibleWorlds {
	hash := game.BuildActionWorldHash(action, possibleWorlds.GetHash())
	response := string(reaction.Response)
	return game.TheCache.GetNextWorlds(hash, response)
}

// ###########################
// ### THE GAME ###
// ###########################
func PlayGames(agreement game.Agreement, motivation game.Motivation, lines []string, limit int) []game.GamePlay {
	game.BUniverse = game.InitUniverse(lines)
	// defer profile.Start().Stop()
	games := []game.GamePlay{}
	for j := uint64(0); j < uint64(limit); j++ {
		perspective := game.NewPerspective(
			motivation,
			agreement,
		)

		btruth := game.BUniverse.GetWorldByIndex(j)
		truth := game.World(btruth.Name)
		gameplay, _ := PlayGame(truth, perspective)
		// fmt.Printf("TRUTH: %s | ", truth)
		// fmt.Printf("TURNS: %d\n", gameplay.TurnCount)
		games = append(games, gameplay)
	}

	return games
}

func PlayGame(truth game.World, perspective game.Perspective) (game.GamePlay, game.World) {
	possibleWorlds := game.BUniverse.PossibleWorlds
	gameplay := game.GamePlay{} // Should this be on perspective?
	reaction := game.Reaction{}
	for true {
		reaction, possibleWorlds = playTurn(gameplay, truth, perspective, possibleWorlds)
		gameplay.AddTurn(reaction)
		if string(reaction.Action) == string(truth) {
			break
		}
	}
	return gameplay, game.World(reaction.Action)
}

func playTurn(
	gameplay game.GamePlay,
	truth game.World,
	perspective game.Perspective,
	possibleWorlds game.PossibleWorlds,
) (
	game.Reaction,
	game.PossibleWorlds,
) {
	predictions := Explore(gameplay, perspective.Agreement, possibleWorlds)
	action := Decide(predictions, perspective.Motivation)
	reaction := Act(gameplay, action, perspective.Agreement.Act, truth)
	possibleWorlds = Observe(action, reaction, possibleWorlds)
	return reaction, possibleWorlds
}
