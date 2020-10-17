package play

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/utils"

	"github.com/chrisconley/go-hangman/test"
)

func TestAct(t *testing.T) {
	truth := game.NewWorld("abc")
	actual := Act(game.InitGamePlay(), game.NewAction("a"), game.BuildTestAct(0.222), truth)
	test.AssertDeepEqual(
		actual,
		game.Reaction{Action: game.NewAction("a"), Response: "a", Reward: 0.222},
		t,
	)
}

func TestObserve(t *testing.T) {
	t.Skip("TODO: Fix test setup")
}

func TestPlayBHangman(t *testing.T) {
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

	for _, line := range lines {
		perspective := game.NewPerspective(
			game.Motivation{Reward: 0.5, Entropy: 0.25, Time: 0.25},
			hangman.Agreement,
		)
		truth := game.NewWorld(line)
		_, actual := PlayGame(truth, perspective)
		test.AssertDeepEqual(actual, truth, t)
	}
}

func TestPlayMastermind(t *testing.T) {
	lines := []string{
		"111",
		"112",
		"121",
		"211",
	}
	game.BUniverse = game.InitUniverse(lines)

	for _, line := range lines {
		perspective := game.NewPerspective(
			game.Motivation{Reward: 0.5, Entropy: 0.25, Time: 0.25},
			mastermind.Agreement,
		)

		truth := game.NewWorld(line)
		_, actual := PlayGame(truth, perspective)
		test.AssertDeepEqual(actual, truth, t)
	}
}
func TestIntegrationPlayKnuthMastermind(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestIntegrationPlayKnuthMastermind in short mode.")
	}
	lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")
	game.BUniverse = game.InitUniverse(lines)

	perspective := game.NewPerspective(
		game.Motivation{Entropy: 0.00, MiniMax: 1.0},
		mastermind.Agreement,
	)

	truth := game.NewWorld("CFCB")
	gameplay, _ := PlayGame(truth, perspective)

	test.AssertDeepEqual(gameplay.Reactions[0].Action, game.NewAction("AABB"), t)
	test.AssertDeepEqual(gameplay.Reactions[0].Response, game.Response("B"), t)

	test.AssertDeepEqual(gameplay.Reactions[1].Action, game.NewAction("ACDD"), t)
	test.AssertDeepEqual(gameplay.Reactions[1].Response, game.Response("W"), t)

	test.AssertDeepEqual(gameplay.Reactions[2].Action, game.NewAction("CEBF"), t)
	test.AssertDeepEqual(gameplay.Reactions[2].Response, game.Response("BWW"), t)

	test.AssertDeepEqual(gameplay.Reactions[3].Action, game.NewAction("ADFB"), t)
	test.AssertDeepEqual(gameplay.Reactions[3].Response, game.Response("BW"), t)
}

func TestIntegrationVerifyMiniMaxMastermind(t *testing.T) {
	t.Skip("Skipping TestIntegrationVerifyMiniMaxMastermind because it's hanging.")
	if testing.Short() {
		t.Skip("Skipping TestIntegrationVerifyMiniMaxMastermind in short mode.")
	}
	lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")

	motivation := game.Motivation{MiniMax: 1.0}
	limit := len(lines)
	games := PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := CalculateStats(games)
	test.AssertDeepEqual(stats.GameCount, 1296, t)
	test.AssertDeepEqual(stats.MaxTurns, 5, t)
	test.AssertDeepEqual(stats.TotalTurns, 5801, t)
	test.AssertAlmostEqual(stats.AverageTurns, 4.4761, 4, t)
}

func TestIntegrationVerifyEntropyMastermind(t *testing.T) {
	t.Skip("Skipping TestIntegrationVerifyEntropyMastermind because it's hanging.")
	if testing.Short() {
		t.Skip("Skipping TestIntegrationVerifyEntropyMastermind in short mode.")
	}
	lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")

	motivation := game.Motivation{Entropy: 1.0}
	limit := len(lines)
	games := PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := CalculateStats(games)
	test.AssertDeepEqual(stats.GameCount, 1296, t)
	test.AssertDeepEqual(stats.MaxTurns, 6, t)
	test.AssertDeepEqual(stats.TotalTurns, 5723, t)
	test.AssertAlmostEqual(stats.AverageTurns, 4.4159, 4, t)
}

func TestIntegrationVerifyEntropyHangman(t *testing.T) {
	t.Skip("Skipping TestIntegrationVerifyEntropyHangman because it's hanging.")
	if testing.Short() {
		t.Skip("Skipping TestIntegrationVerifyEntropyHangman in short mode.")
	}
	lines := utils.ReadLinesFromFilePath("../../fixtures/nine-letter-words.txt")

	motivation := game.Motivation{Entropy: 1.0}
	limit := len(lines)
	games := PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := CalculateStats(games)
	test.AssertDeepEqual(stats.GameCount, 25011, t)
	test.AssertDeepEqual(stats.MaxTurns, 13, t)
	test.AssertDeepEqual(stats.TotalTurns, 146356, t)
	test.AssertAlmostEqual(stats.AverageTurns, 5.8529, 4, t)
}

func TestIntegrationVerifyMiniMaxHangman(t *testing.T) {
	t.Skip("Skipping TestIntegrationVerifyMiniMaxHangman because it's hanging.")
	if testing.Short() {
		t.Skip("Skipping TestIntegrationVerifyMiniMaxHangman in short mode.")
	}
	lines := utils.ReadLinesFromFilePath("../../fixtures/nine-letter-words.txt")

	motivation := game.Motivation{MiniMax: 1.0}
	limit := len(lines)
	games := PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := CalculateStats(games)
	test.AssertDeepEqual(stats.GameCount, 25011, t)
	test.AssertDeepEqual(stats.MaxTurns, 12, t)
	test.AssertDeepEqual(stats.TotalTurns, 149552, t)
	test.AssertAlmostEqual(stats.AverageTurns, 5.9794, 4, t)
}

func TestPlayGames(t *testing.T) {
	lines := utils.ReadLinesFromFilePath("../../fixtures/nine-letter-words.txt")
	lines = lines[:100]

	motivation := game.Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
	games := PlayGames(mastermind.Agreement, motivation, lines, 3)
	test.AssertDeepEqual(games[0].TurnCount, 2, t)
	test.AssertDeepEqual(games[1].TurnCount, 3, t)
	test.AssertDeepEqual(games[2].TurnCount, 3, t)
}

func TestBuildTrialsTwoDimensions(t *testing.T) {
	trials := BuildTrials(0.0, 1.0, 0.25, 2)
	test.AssertDeepEqual(len(trials), 5, t)
	test.AssertDeepEqual(trials[0], []float64{0.0, 1.0}, t)
	test.AssertDeepEqual(trials[1], []float64{0.25, 0.75}, t)
	test.AssertDeepEqual(trials[2], []float64{0.5, 0.5}, t)
	test.AssertDeepEqual(trials[3], []float64{0.75, 0.25}, t)
	test.AssertDeepEqual(trials[4], []float64{1.0, 0.0}, t)
}
func TestBuildTrialsThreeDimensions(t *testing.T) {
	trials := BuildTrials(0.0, 1.0, 0.25, 3)
	test.AssertDeepEqual(len(trials), 15, t)
	test.AssertDeepEqual(trials[0], []float64{0.0, 0.0, 1.0}, t)
	test.AssertDeepEqual(trials[1], []float64{0.0, 0.25, 0.75}, t)
	test.AssertDeepEqual(trials[2], []float64{0.0, 0.5, 0.5}, t)
	test.AssertDeepEqual(trials[3], []float64{0.0, 0.75, 0.25}, t)
	test.AssertDeepEqual(trials[4], []float64{0.0, 1.0, 0.0}, t)

	test.AssertDeepEqual(trials[5], []float64{0.25, 0.0, 0.75}, t)
	test.AssertDeepEqual(trials[6], []float64{0.25, 0.25, 0.5}, t)
	test.AssertDeepEqual(trials[7], []float64{0.25, 0.5, 0.25}, t)
	test.AssertDeepEqual(trials[8], []float64{0.25, 0.75, 0.0}, t)

	test.AssertDeepEqual(trials[9], []float64{0.5, 0.0, 0.5}, t)
	test.AssertDeepEqual(trials[10], []float64{0.5, 0.25, 0.25}, t)
	test.AssertDeepEqual(trials[11], []float64{0.5, 0.5, 0.0}, t)

	test.AssertDeepEqual(trials[12], []float64{0.75, 0.0, 0.25}, t)
	test.AssertDeepEqual(trials[13], []float64{0.75, 0.25, 0.0}, t)

	test.AssertDeepEqual(trials[14], []float64{1.0, 0.0, 0.0}, t)
}
func TestBuildTrialsFourDimensions(t *testing.T) {
	trials := BuildTrials(0.0, 1.0, 0.25, 4)
	test.AssertDeepEqual(len(trials), 35, t)
	test.AssertDeepEqual(trials[0], []float64{0.0, 0.0, 0.0, 1.0}, t)
	test.AssertDeepEqual(trials[14], []float64{0.0, 1.0, 0.0, 0.0}, t)
	test.AssertDeepEqual(trials[15], []float64{0.25, 0.0, 0.0, 0.75}, t)
	test.AssertDeepEqual(trials[24], []float64{0.25, 0.75, 0.0, 0.0}, t)
	test.AssertDeepEqual(trials[25], []float64{0.5, 0.0, 0.0, 0.5}, t)
	test.AssertDeepEqual(trials[30], []float64{0.5, 0.5, 0.0, 0.0}, t)
	test.AssertDeepEqual(trials[31], []float64{0.75, 0.0, 0.0, 0.25}, t)
	test.AssertDeepEqual(trials[33], []float64{0.75, 0.25, 0.0, 0.0}, t)
	test.AssertDeepEqual(trials[34], []float64{1.0, 0.0, 0.0, 0.0}, t)
}

func TestBuildTrialsFourDimensionSmallerStep(t *testing.T) {
	trials := BuildTrials(0.0, 1.0, 0.2, 4)
	test.AssertDeepEqual(len(trials), 56, t)

	trials = BuildTrials(0.0, 1.0, 0.1, 3)
	test.AssertDeepEqual(len(trials), 66, t)

	trials = BuildTrials(0.0, 1.0, 0.05, 3)
	test.AssertDeepEqual(len(trials), 231, t)

	trials = BuildTrials(0.0, 1.0, 0.01, 3)
	test.AssertDeepEqual(len(trials), 5151, t)
}
