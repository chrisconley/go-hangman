package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	// "io"
	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"

	// "github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/play"
)

func main() {
	examplePtr := flag.String("example", "defaultValue", " Help text.")
	flag.Parse()
	fmt.Printf("\n EXAMPLE: %s\n", *examplePtr)

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: fortune | gocowsay")
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("ERROR")
	}

	worlds := []game.World{}
	for j := 0; j < len(lines); j++ {
		worlds = append(worlds, game.World(game.World(lines[j])))
	}
	possibleWorlds := game.BuildPossibleWorlds(worlds)
	// fmt.Printf("\nWORLDS: %+v", worlds)

	turns := []int{}
	for j := 0; j < 100; j++ {
		perspective := &game.Perspective{}
		perspective.Knowledge = game.Knowledge{}
		perspective.PossibleWorlds = possibleWorlds
		perspective.AvailableActions = game.BuildAvailableActions(possibleWorlds)
		// fmt.Printf("\nAVAILABLE ACTIONS: %+v", perspective.AvailableActions)
		perspective.Agreement = hangman.HangmanAgreement
		perspective.ActionFilter = hangman.ActionFilter
		// perspective.Agreement = mastermind.Agreement
		perspective.Motivation = game.Motivation{0.0, 1.0, 0.0}

		truth := game.World(string(lines[j]))
		fmt.Printf("TRUTH: %s | ", truth)
		// memory := game.InitMemory(&game.NaiveMemory{})
		memory := game.InitMemory(&game.NoMemory{})
		result := play.PlayGame(truth, perspective, memory)
		fmt.Printf("RESULT: %s | ", result)
		fmt.Printf("TURNS: %d\n", len(perspective.Knowledge))
		turns = append(turns, len(perspective.Knowledge))
	}
	sum := 0
	for i := 0; i < len(turns); i++ {
		sum = sum + turns[i]
	}
	avg := float64(sum) / float64(len(turns))
	fmt.Printf("AVG TURNS: %f\n", avg)
}
