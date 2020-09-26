package main

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/chrisconley/go-hangman/test"
)

func TestIntegration(t *testing.T) {
	lines := utils.ReadLinesFromFilePath("../fixtures/nine-letter-words.txt")
	lines = lines[:100]

	motivation := game.Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
	turns := play.PlayGames(mastermind.Agreement, motivation, lines, 5)
	stats := play.CalculateStats(turns)

	test.AssertAlmostEqual(stats.AverageTurns, 1.8, 4, t)
}
