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
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
)

func main() {
	// ballast := make([]byte, 10<<28)
	limit := flag.Int("limit", 10, " Help text.")
	flag.Parse()

	lines, err := utils.ReadLinesFromStdin()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	game.InitUniverse(lines)

	worlds := game.BuildPossibleWorlds(lines)
	// defer profile.Start().Stop()

	turns := []int{}
	for j := 0; j < *limit; j++ {
		perspective := game.NewPerspective(
			game.Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0},
			mastermind.Agreement,
		)

		truth := worlds.Worlds[j]
		turnCount, _ := play.PlayGame(truth, perspective)
		// fmt.Printf("TRUTH: %s | ", truth)
		// fmt.Printf("TURNS: %d\n", turnCount)
		turns = append(turns, turnCount)
	}
	sum := 0
	for i := 0; i < len(turns); i++ {
		sum = sum + turns[i]
	}
	avg := float64(sum) / float64(len(turns))
	fmt.Printf("AVG TURNS: %.4f | (sum: %d, games: %d)\n", avg, sum, len(turns))
	// count := 0
	// hashes := make(map[string]bool)
	// for _, worlds := range game.BBPossibleWorldsStore {
	// 	for _, world := range worlds {
	// 		count++
	// 		hashes[string(world.GetHash())] = true
	// 		if count <= 3 {
	// 			fmt.Printf("hash: %+v\n", string(world.GetHash()))
	// 		}
	// 	}
	// }
	// fmt.Printf("Universe store key size: %d\n", len(game.BBPossibleWorldsStore))
	// fmt.Printf("Universe store value size: %d\n", count)
	// fmt.Printf("Unique worlds size: %d\n", len(hashes))
	// fmt.Printf("Ballast %v\n", ballast[0])
}
