package main

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"
	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/chrisconley/go-hangman/test"
)

func TestAverageGuesses(t *testing.T) {
	lines := utils.ReadLinesFromFilePath("../fixtures/nine-letter-words.txt")
	lines = lines[:1000]

	possibleWorlds := game.BuildPossibleWorlds(lines)
	game.InitUniverse(lines)
	turns := []int{}
	for j := 0; j < 5; j++ {
		perspective := game.NewPerspective(
			game.Motivation{Reward: 0.5, Entropy: 0.25, Time: 0.25},
			hangman.Agreement,
		)

		truth := possibleWorlds.Worlds[j]
		turnCount, _ := play.PlayGame(truth, perspective)
		turns = append(turns, turnCount)
	}
	sum := 0
	for i := 0; i < len(turns); i++ {
		sum = sum + turns[i]
	}
	avg := float64(sum) / float64(len(turns))
	test.AssertAlmostEqual(avg, 2.6000, 4, t) //TODO This should really be 4.2

}
