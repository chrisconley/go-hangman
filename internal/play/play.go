package play

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/coocood/freecache"
	"github.com/segmentio/fasthash/fnv1a"
)

func logPredictions(predictions game.Predictions) {
	if os.Getenv("DEBUGPREDICTIONS") != "" {
		for _, prediction := range predictions.PredictionsMap.Values {
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

func logProgress(gameIndex uint64, truth game.World) {
	if gameIndex%100 == 0 && os.Getenv("LOGPROGRESS") != "" {
		fmt.Printf("%s | ", time.Now().Format("2006.01.02 15:04:05"))
		fmt.Printf("GAME #: %d | ", gameIndex)
		fmt.Printf("TRUTH: %s \n", truth)
	}
}

// ###########################
// ### THE DECISION PROCESS ###
// ###########################
func Explore(gameplay game.GamePlay, agreement game.Agreement, possibleWorlds game.PossibleWorlds) (predictions game.Predictions) {
	gameHash := gameplay.Hash
	predictions, err := game.PredictionCache.Get(gameHash)
	if err != nil {
		actions := agreement.GetActions(
			gameplay,
			possibleWorlds,
			game.BUniverse.Actions,
		)
		predictions = game.ExploreReactions(gameplay, possibleWorlds, actions, agreement.Act)
		game.PredictionCache.Set(gameHash, predictions)
	}
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

var TheTestCache = freecache.NewCache(1000 * 1024 * 1024)

func Observe(gameplay game.GamePlay, reaction game.Reaction, act game.Act, possibleWorlds game.PossibleWorlds) game.PossibleWorlds {
	hash := fnv1a.AddUint64(reaction.Response.Hash(), reaction.Action.Hash())
	hash = fnv1a.AddUint64(hash, gameplay.Hash)
	hash = fnv1a.AddUint64(hash, possibleWorlds.SetHash())
	npw, err := game.ObserveCache.Get(hash)
	if err != nil {
		npw = game.NewPossibleWorlds(game.BUniverse.Size)
		for _, worldIndex := range possibleWorlds.Indices() {
			world := game.BUniverse.GetWorldByIndex(worldIndex)
			testReaction := act(gameplay, world, reaction.Action)
			if testReaction.Response.Matches(reaction.Response) {
				npw.Add(worldIndex)
			}
		}
		game.ObserveCache.Set(hash, npw)
	}
	return npw
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

		truth := game.BUniverse.GetWorldByIndex(j)
		gameplay, _ := PlayGame(truth, perspective)
		logProgress(j+1, truth)
		games = append(games, gameplay)
	}

	return games
}

func PlayGame(truth game.World, perspective game.Perspective) (game.GamePlay, game.World) {
	possibleWorlds := game.BUniverse.PossibleWorlds
	gameplay := game.InitGamePlay(possibleWorlds.GetHash()) // Should this be on perspective?
	reaction := game.Reaction{}
	for true {
		reaction, possibleWorlds = playTurn(gameplay, truth, perspective, possibleWorlds)
		gameplay.AddTurn(reaction)
		if reaction.Action.Matches(truth) {
			break
		}
	}
	return gameplay, truth
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
	possibleWorlds = Observe(gameplay, reaction, perspective.Agreement.Act, possibleWorlds)
	return reaction, possibleWorlds
}

func BuildTrials(min float64, max float64, stepSize float64, numDimensions int) [][]float64 {
	values := []float64{}
	currentValue := min
	for true {
		values = append(values, currentValue)
		if math.Abs(currentValue-max) <= 1e-9 {
			break
		} else if currentValue-max > 1e-9 {
			fmt.Printf("min %f, max %f, stepSize %f, currentValue %f", min, max, stepSize, currentValue)
			panic("stepSize is not equally divisible")
		}
		currentValue += stepSize
	}

	trials := [][]float64{}
	if numDimensions == 2 {
		for _, value := range values {
			trials = append(trials, []float64{value, max - value})
		}
		return trials
	}

	for _, value := range values {
		newTrials := BuildTrials(min, max-value, stepSize, numDimensions-1)
		for _, newTrial := range newTrials {
			trials = append(trials, append([]float64{value}, newTrial...))
		}
	}
	return trials
}
