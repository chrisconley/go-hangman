package main 

import (
	"flag"
	"fmt"
	"os"
	"bufio"
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

	originalWorlds := game.PossibleWorlds{}
	originalActions := game.PossibleActions{}
	for j := 0; j < len(lines); j++ {
		originalWorlds = append(originalWorlds, game.PossibleWorld{game.World(lines[j]), 0.5})
		// originalActions = append(originalActions, game.Action(lines[j]))
	}
	originalActions = game.PossibleActions{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}


	turns := []int{}
	for j := 0; j < 100; j++ {
		perspective := &game.Perspective{}
		perspective.Knowledge = game.Knowledge{}
		perspective.PossibleWorlds = originalWorlds
		perspective.PossibleActions = originalActions
		perspective.Agreement = hangman.HangmanAgreement
		// perspective.Agreement = mastermind.Agreement
		perspective.Motivation = game.Motivation{0.0, 1.0, 0.0}

		truth := game.World(string(lines[j]))
		fmt.Printf("TRUTH: %s | ", truth)
		result := play.PlayGame(truth, perspective)
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