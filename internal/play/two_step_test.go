package play

import (
	"fmt"
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/utils"
)

func TestTwoStep(t *testing.T) {
	t.Skip("TODO: Either fix TwoStep or remove")
	lines := []string{
		"AA", "AB", "AC", "AD",
		"BA", "BB", "BC", "BD",
		"CA", "CB", "CC", "CD",
		"DA", "DB", "DC", "DD",
	}
	game.BUniverse = game.InitUniverse(lines)
	startingGamePlay := game.InitGamePlay(game.BUniverse.PossibleWorlds.GetHash())
	firstStepPredictions := Explore(startingGamePlay, mastermind.Agreement, game.BUniverse.PossibleWorlds)
	for _, prediction := range firstStepPredictions.PredictionsMap.Values {
		fmt.Printf("%+v\n", prediction)
	}
	predictions := GetTwoStepPredictions(
		startingGamePlay,
		mastermind.Agreement,
		game.BUniverse.PossibleWorlds,
	)
	for _, prediction := range predictions.PredictionsMap.Values {
		fmt.Printf("hi %+v\n", prediction)
	}
	// test.AssertDeepEqual(1, 2, t)
}

func TestTwoStepLarge(t *testing.T) {
	t.Skip("TODO: Either fix TwoStep or remove")
	// lines := utils.ReadLinesFromFilePath("../../fixtures/mastermind-ABCD-3.txt")
	lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")
	game.BUniverse = game.InitUniverse(lines)
	startingGamePlay := game.InitGamePlay(game.BUniverse.PossibleWorlds.GetHash())
	firstStepPredictions := Explore(startingGamePlay, mastermind.Agreement, game.BUniverse.PossibleWorlds)
	for _, prediction := range firstStepPredictions.PredictionsMap.Values {
		fmt.Printf("%+v\n", prediction)
	}
	predictions := GetTwoStepPredictions(
		startingGamePlay,
		mastermind.Agreement,
		game.BUniverse.PossibleWorlds,
	)
	for _, prediction := range predictions.PredictionsMap.Values {
		fmt.Printf("hi %+v\n", prediction)
	}
	// test.AssertDeepEqual(1, 2, t)
}

func TestTwoStepLast(t *testing.T) {
	t.Skip("TODO: Either fix TwoStep or remove")
	lines := []string{
		"aardvarks",
		"aasvogels",
		"abamperes",
		"abandoned",
		"abandoner",
		"abasement",
		"abashment",
		"abatement",
		"abattises",
		"abattoirs",
	}
	game.BUniverse = game.InitUniverse(lines)
	startingGamePlay := game.InitGamePlay(game.BUniverse.PossibleWorlds.GetHash())
	possibleWorlds := game.NewPossibleWorlds(game.BUniverse.Size)
	possibleWorlds.Add(1)
	predictions := GetTwoStepPredictions(
		startingGamePlay,
		hangman.Agreement,
		possibleWorlds,
	)
	for _, prediction := range predictions.PredictionsMap.Values {
		fmt.Printf("hi %+v\n", prediction)
	}
	// test.AssertDeepEqual(1, 2, t)
}
