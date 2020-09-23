// Example usage
// cat ./fixtures/nine-letter-words.txt | go run cmd/player.go

// Testing against python hangman
// head -n 1000 ./fixtures/nine-letter-words.txt | go run cmd/player.go -limit 100 > /tmp/go-hangman.tx
// head -n 1000 ./build/splits/9 | python3 ./games/play.py - --game hangman --limit 100 --strategy "weighted_sum|reward:0;info:100;speed:0" > /tmp/hangman.txt
// diff -u /tmp/hangman.txt /tmp/go-hangman.tx | ydiff -s

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"
	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/pkg/profile"
)

func main() {
	// ballast := make([]byte, 10<<30)
	limit := flag.Int("limit", 10, " Help text.")
	flag.Parse()

	lines, err := utils.ReadLinesFromStdin()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	game.InitUniverse(lines)

	worlds := game.BuildPossibleWorlds(lines)
	defer profile.Start().Stop()

	turns := []int{}
	for j := 0; j < *limit; j++ {
		perspective := game.NewPerspective(
			game.Motivation{Reward: 0.1, Entropy: 0.9, Time: 0.0},
			hangman.HangmanAgreement,
			hangman.ActionFilter,
			hangman.HangmanBAgreement,
		)

		truth := worlds.Worlds[j]
		turnCount, _ := play.PlayGame(truth, perspective)
		fmt.Printf("TRUTH: %s | ", truth)
		fmt.Printf("TURNS: %d\n", turnCount)
		turns = append(turns, turnCount)
	}
	sum := 0
	for i := 0; i < len(turns); i++ {
		sum = sum + turns[i]
	}
	avg := float64(sum) / float64(len(turns))
	fmt.Printf("AVG TURNS: %.4f | (sum: %d, games: %d)\n", avg, sum, len(turns))
	// fmt.Printf("Ballast %v", ballast[0])
}
