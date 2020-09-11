// Example usage
// cat ./fixtures/nine-letter-words.txt | go run cmd/player.go

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"

	// "io"
	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"

	// "github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/play"
)

func readLinesFromStdin() ([]string, error) {
	var lines []string
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		return lines, errors.New("The command is intended to work with pipes: cat ./fixtures/nine-letter-words.txt | go run cmd/player.go")
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("ERROR")
	}

	return lines, nil
}

func main() {
	examplePtr := flag.String("example", "defaultValue", " Help text.")
	flag.Parse()
	fmt.Printf("\n EXAMPLE: %s\n", *examplePtr)

	lines, err := readLinesFromStdin()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	possibleWorlds := game.BuildPossibleWorlds(lines)

	turns := []int{}
	for j := 0; j < 4; j++ {
		fmt.Printf("possibleWorlds length: %d | ", possibleWorlds.Len())
		motivation := game.Motivation{0.0, 1.0, 0.0}
		perspective := game.NewPerspective(possibleWorlds, motivation)
		perspective.Agreement = hangman.HangmanAgreement
		perspective.ActionFilter = hangman.ActionFilter

		truth := possibleWorlds.Worlds[j]
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
