package main

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/chrisconley/go-hangman/test"
)

func TestIntegration(t *testing.T) {
	lines := utils.ReadLinesFromFilePath("../fixtures/nine-letter-words.txt")
	lines = lines[:100]

	turns := play.PlayGames(lines, 5)
	stats := play.CalculateStats(turns)

	test.AssertAlmostEqual(stats.AverageTurns, 1.8, 4, t)
}
