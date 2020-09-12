// Example usage
// cat ./fixtures/nine-letter-words.txt | go run cmd/player.go

// Testing against python hangman
// head -n 1000 ./fixtures/nine-letter-words.txt | go run cmd/player.go -limit 100 > /tmp/go-hangman.tx
// head -n 1000 ./build/splits/9 | python3 ./games/play.py - --game hangman --limit 100 --strategy "weighted_sum|reward:0;info:100;speed:0" > /tmp/hangman.txt
// diff -u /tmp/hangman.txt /tmp/go-hangman.tx | ydiff -s

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
	"github.com/chrisconley/go-hangman/internal/play"
	// "github.com/chrisconley/go-hangman/internal/mastermind"
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
	limit := flag.Int("limit", 10, " Help text.")
	flag.Parse()

	lines, err := readLinesFromStdin()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	possibleWorlds := game.BuildPossibleWorlds(lines)

	turns := []int{}
	for j := 0; j < *limit; j++ {
		motivation := game.Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
		perspective := game.NewPerspective(possibleWorlds, motivation)
		perspective.Agreement = hangman.HangmanAgreement
		perspective.ActionFilter = hangman.ActionFilter

		truth := possibleWorlds.Worlds[j]
		// memory := game.InitMemory(&game.NaiveMemory{})
		memory := game.InitMemory(&game.NoMemory{})
		_ = play.PlayGame(truth, perspective, memory)
		// fmt.Printf("RESULT: %s | ", result)
		fmt.Printf("TRUTH: %s | ", truth)
		fmt.Printf("TURNS: %d\n", len(perspective.Knowledge))
		turns = append(turns, len(perspective.Knowledge))
	}
	sum := 0
	for i := 0; i < len(turns); i++ {
		sum = sum + turns[i]
	}
	avg := float64(sum) / float64(len(turns))
	fmt.Printf("AVG TURNS: %.4f\n", avg)
}
