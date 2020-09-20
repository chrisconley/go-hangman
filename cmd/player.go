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

	// "io"
	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"
	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
	// "github.com/chrisconley/go-hangman/internal/mastermind"
)

func main() {
	limit := flag.Int("limit", 10, " Help text.")
	flag.Parse()

	lines, err := utils.ReadLinesFromStdin()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	universeLength := uint64(len(lines))
	game.BUniverse = make([]game.BWorld, universeLength)
	fmt.Printf("UNVERSE LEN: %d", universeLength)
	pw := game.NewPossibleWorlds(universeLength)
	for index, name := range lines {
		game.BUniverse[index] = game.BWorld{Index: uint64(index), Name: name}
		pw.Bits.SetBit(uint64(index))
	}

	// defer profile.Start().Stop()
	worlds := game.BuildPossibleWorlds(lines)

	turns := []int{}
	for j := 0; j < *limit; j++ {
		perspective := &game.Perspective{}
		perspective.Knowledge = game.Knowledge{}
		// perspective.PossibleWorlds = worlds
		perspective.BPossibleWorlds = pw
		perspective.AvailableActions = game.BuildAvailableActions(worlds)
		perspective.Agreement = hangman.HangmanAgreement
		perspective.ActionFilter = hangman.ActionFilter
		perspective.Motivation = game.Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}

		truth := worlds.Worlds[j]
		fmt.Printf("TRUTH: %s | ", truth)
		// memory := game.InitMemory(&game.NaiveMemory{})
		memory := game.InitMemory(&game.NoMemory{})
		_ = play.PlayGame(truth, perspective, memory)
		// fmt.Printf("RESULT: %s | ", result)
		fmt.Printf("TURNS: %d\n", len(perspective.Knowledge))
		turns = append(turns, len(perspective.Knowledge))
	}
	sum := 0
	for i := 0; i < len(turns); i++ {
		sum = sum + turns[i]
	}
	avg := float64(sum) / float64(len(turns))
	fmt.Printf("AVG TURNS: %.4f | (sum: %d, games: %d)\n", avg, sum, len(turns))
}
