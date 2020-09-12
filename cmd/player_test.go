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
	lines := utils.LoadLinesFromFilePath("../fixtures/nine-letter-words.txt")
	possibleWorlds := game.BuildPossibleWorlds(lines[:1000])
	turns := []int{}
	for j := 0; j < 5; j++ {
		motivation := game.Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
		perspective := game.NewPerspective(possibleWorlds, motivation)
		perspective.Agreement = hangman.HangmanAgreement
		perspective.ActionFilter = hangman.ActionFilter

		truth := possibleWorlds.Worlds[j]
		memory := game.InitMemory(&game.NoMemory{})
		_ = play.PlayGame(truth, perspective, memory)
		turns = append(turns, len(perspective.Knowledge))
	}
	sum := 0
	for i := 0; i < len(turns); i++ {
		sum = sum + turns[i]
	}
	avg := float64(sum) / float64(len(turns))
	test.AssertAlmostEqual(avg, 4.2000, 4, t)

}
