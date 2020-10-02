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
	truth := game.World("abc")
	actual := Act(game.Action("a"), game.BuildTestAct(0.222), truth)
	test.AssertDeepEqual(
		actual,
		game.Reaction{Action: "a", Response: "a", Reward: 0.222},
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
		truth := game.World(line)
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

		truth := game.World(line)
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

	truth := game.World("CFCB")
	reactions, _ := PlayGame(truth, perspective)

	test.AssertDeepEqual(reactions.Turns[0].Action, game.Action("AABB"), t)
	test.AssertDeepEqual(reactions.Turns[0].Response, game.Response("B"), t)

	test.AssertDeepEqual(reactions.Turns[1].Action, game.Action("ACDD"), t)
	test.AssertDeepEqual(reactions.Turns[1].Response, game.Response("W"), t)

	test.AssertDeepEqual(reactions.Turns[2].Action, game.Action("CEBF"), t)
	test.AssertDeepEqual(reactions.Turns[2].Response, game.Response("BWW"), t)

	test.AssertDeepEqual(reactions.Turns[3].Action, game.Action("ADFB"), t)
	test.AssertDeepEqual(reactions.Turns[3].Response, game.Response("BW"), t)
}

func TestIntegrationVerifyMiniMaxMastermind(t *testing.T) {
	t.Skip("Skipping TestIntegrationVerifyMiniMaxMastermind because it's hanging.")
	if testing.Short() {
		t.Skip("Skipping TestIntegrationVerifyMiniMaxMastermind in short mode.")
	}
	lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")

	motivation := game.Motivation{MiniMax: 1.0}
	limit := len(lines)
	turns := PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := CalculateStats(turns)
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
	turns := PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := CalculateStats(turns)
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
	turns := PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := CalculateStats(turns)
	test.AssertDeepEqual(stats.GameCount, 25011, t)
	test.AssertDeepEqual(stats.MaxTurns, 13, t)
	test.AssertDeepEqual(stats.TotalTurns, 146387, t)
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
	turns := PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := CalculateStats(turns)
	test.AssertDeepEqual(stats.GameCount, 25011, t)
	test.AssertDeepEqual(stats.MaxTurns, 12, t)
	test.AssertDeepEqual(stats.TotalTurns, 149552, t)
	test.AssertAlmostEqual(stats.AverageTurns, 5.9794, 4, t)
}

func TestPlayGames(t *testing.T) {
	lines := utils.ReadLinesFromFilePath("../../fixtures/nine-letter-words.txt")
	lines = lines[:100]

	motivation := game.Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
	turns := PlayGames(mastermind.Agreement, motivation, lines, 3)
	test.AssertDeepEqual(turns[0].TurnCount, 2, t)
	test.AssertDeepEqual(turns[1].TurnCount, 3, t)
	test.AssertDeepEqual(turns[2].TurnCount, 3, t)
}
