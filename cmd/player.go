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

func PlayGames(lines []string, limit int) []int {
	game.InitUniverse(lines)
	// defer profile.Start().Stop()
	turns := []int{}
	for j := 0; j < limit; j++ {
		perspective := game.NewPerspective(
			game.Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0},
			mastermind.Agreement,
		)

		btruth := game.BUniverse.Worlds[j]
		truth := game.World(btruth.Name)
		turnCount, _ := play.PlayGame(truth, perspective)
		// fmt.Printf("TRUTH: %s | ", truth)
		// fmt.Printf("TURNS: %d\n", turnCount)
		turns = append(turns, turnCount)
	}

	return turns
}

type Stats struct {
	TotalTurns   int
	AverageTurns float64
	MaxTurns     int
	GameCount    int
}

func CalculateStats(turnCounts []int) Stats {
	stats := Stats{}
	stats.TotalTurns = 0
	stats.MaxTurns = 0
	stats.GameCount = len(turnCounts)
	for _, turnCount := range turnCounts {
		stats.TotalTurns = stats.TotalTurns + turnCount
		if turnCount > stats.MaxTurns {
			stats.MaxTurns = turnCount
		}
	}

	stats.AverageTurns = float64(stats.TotalTurns) / float64(stats.GameCount)
	fmt.Printf(
		"AVG TURNS: %.4f | (sum: %d, games: %d)\n",
		stats.AverageTurns,
		stats.TotalTurns,
		stats.GameCount,
	)
	return stats
}

func main() {
	lines, err := utils.ReadLinesFromStdin()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	limit := flag.Int("limit", len(lines), "How many games to play. Defaults to all games piped in from stdin.")
	flag.Parse()

	turns := PlayGames(lines, *limit)
	_ = CalculateStats(turns)
}
