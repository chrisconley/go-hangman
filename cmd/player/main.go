// Example usage
// cat ./fixtures/nine-letter-words.txt | go run cmd/player/main.go

// Testing against python hangman
// head -n 1000 ./fixtures/nine-letter-words.txt | go run cmd/player/main.go -limit 100 > /tmp/go-hangman.tx
// head -n 1000 ./build/splits/9 | python3 ./games/play.py - --game hangman --limit 100 --strategy "weighted_sum|reward:0;info:100;speed:0" > /tmp/hangman.txt
// diff -u /tmp/hangman.txt /tmp/go-hangman.tx | ydiff -s

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/pkg/profile"

	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
)

func main() {
	lines, err := utils.ReadLinesFromStdin()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	limit := flag.Int("limit", len(lines), "How many games to play. Defaults to all games piped in from stdin.")
	gameString := flag.String("game", "none", "What game to play. Valid games: hangman, mastermind. Required")
	reward := flag.Float64("motivation.reward", 0.0, "How much weight to put into Reward.")
	entropy := flag.Float64("motivation.entropy", 0.0, "How much weight to put into Entropy.")
	minimax := flag.Float64("motivation.minimax", 0.0, "How much weight to put into MiniMax.")
	profiling := flag.String("profiling", "off", "Set 'mem' or 'cpu' profiling. Defaults to off.")
	profilingProfilePath := flag.String("profiling.profilePath", "/tmp/profiling", "Path where profile should be written")
	flag.Parse()

	if *reward+*entropy+*minimax != 1.0 {
		fmt.Printf("%f %f %f", *reward, *entropy, *minimax)
		panic("Reward + Entropy + Minimax must equal 1.0. All params default to 0.0")
	}

	agreement := hangman.Agreement
	if *gameString == "hangman" {
		agreement = hangman.Agreement
	} else if *gameString == "mastermind" {
		agreement = mastermind.Agreement
	} else {
		panic("A game must be specified.")
	}

	if *profiling == "cpu" {
		defer profile.Start(
			profile.ProfilePath(*profilingProfilePath),
		).Stop()

	} else if *profiling == "mem" {
		defer profile.Start(
			profile.MemProfile,
			profile.MemProfileRate(4096*10),
			profile.ProfilePath(*profilingProfilePath),
		).Stop()
	}

	motivation := game.Motivation{Reward: *reward, Entropy: *entropy, MiniMax: *minimax}

	games := play.PlayGames(agreement, motivation, lines, *limit)
	stats := play.CalculateStats(games)

	b, err := json.Marshal(stats)
	fmt.Printf("GameStats: %s\n", string(b))
	fmt.Printf("PredictionCache stats: %+v\n", game.PredictionCache.Stats())
	fmt.Printf("PredictionCache size: %+v\n", game.PredictionCache.Size())
	fmt.Printf("ReactionCache stats: %+v\n", game.ReactionCache.Stats())
	fmt.Printf("ReactionCache size: %+v\n", game.ReactionCache.Size())
}
