package main

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/chrisconley/go-hangman/test"
)

func TestIntegration(t *testing.T) {
	lines := utils.ReadLinesFromFilePath("../fixtures/nine-letter-words.txt")
	lines = lines[:100]

	turns := PlayGames(lines, 5)
	stats := CalculateStats(turns)

	test.AssertAlmostEqual(stats.AverageTurns, 1.8, 4, t)
}

func TestPlayGames(t *testing.T) {
	lines := utils.ReadLinesFromFilePath("../fixtures/nine-letter-words.txt")
	lines = lines[:100]

	game.InitUniverse(lines)
	turns := PlayGames(lines, 3)
	test.AssertDeepEqual(turns, []int{1, 2, 2}, t)
}

func TestCalculateStats(t *testing.T) {
	turns := []int{1, 2, 3, 5, 11, 7}
	stats := CalculateStats(turns)

	test.AssertDeepEqual(stats.GameCount, 6, t)
	test.AssertDeepEqual(stats.TotalTurns, 29, t)
	test.AssertAlmostEqual(stats.AverageTurns, 4.8333, 4, t)
	test.AssertDeepEqual(stats.MaxTurns, 11, t)
}
