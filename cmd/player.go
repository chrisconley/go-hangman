// Example usage
// cat ./fixtures/nine-letter-words.txt | go run cmd/player.go

// Testing against python hangman
// head -n 1000 ./fixtures/nine-letter-words.txt | go run cmd/player.go -limit 100 > /tmp/go-hangman.tx
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
	motivation := game.Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
	flag.Parse()

	agreement := hangman.Agreement
	if *gameString == "hangman" {
		agreement = hangman.Agreement
	} else if *gameString == "mastermind" {
		agreement = mastermind.Agreement
	} else {
		panic("A game must be specified.")
	}

	turns := play.PlayGames(agreement, motivation, lines, *limit)
	stats := play.CalculateStats(turns)

	b, err := json.Marshal(stats)
	fmt.Printf("GameStats: %s", string(b))
}
