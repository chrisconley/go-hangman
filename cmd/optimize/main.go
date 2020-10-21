package main

import (
	"encoding/json"
	"fmt"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/pkg/profile"
)

type optimizationResults struct {
	Motivation game.Motivation
	GameStats  play.Stats
}

func TryTry(motivation game.Motivation, lines []string) play.Stats {
	// limit := 1000
	limit := len(lines)
	games := play.PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := play.CalculateStats(games)

	js, err := json.Marshal(optimizationResults{Motivation: motivation, GameStats: stats})
	if err != nil {
		fmt.Println(js)
	}
	fmt.Println(string(js))

	return stats
}

func main() {
	// lines := utils.ReadLinesFromFilePath("./fixtures/mastermind-ABCDEFG-4.txt")
	lines := utils.ReadLinesFromFilePath("./fixtures/commercial-mastermind.txt")
	// lines := utils.ReadLinesFromFilePath("./fixtures/nine-letter-words.txt")
	// lines = lines[:10000]
	game.BUniverse = game.InitUniverse(lines)
	trials := play.BuildTrials(0.0, 1.0, 0.5, 4)
	fmt.Println("Number of trials: ", len(trials))
	defer profile.Start(
		// profile.memprofile,
		// profile.memprofilerate(4096*10),
		profile.ProfilePath("/tmp/mem.pprof"),
	).Stop()

	minTotalTurnsMotivation := game.Motivation{}
	minTotalTurns := 1000000

	minTotalPainMotivation := game.Motivation{}
	minTotalPain := 1000000.0
	for _, trial := range trials[:3] {
		for index, value := range trial {
			if value < 0.0 {
				trial[index] = 0.0
			}
		}
		motivation := game.Motivation{
			Reward:    trial[3],
			Entropy:   trial[0],
			MiniMax:   trial[1],
			MostParts: trial[2],
		}
		stats := TryTry(motivation, lines)

		if stats.TotalTurns < minTotalTurns {
			minTotalTurns = stats.TotalTurns
			minTotalTurnsMotivation = motivation
		}
		if stats.TotalPain < minTotalPain {
			minTotalPain = stats.TotalPain
			minTotalPainMotivation = motivation
		}
	}

	fmt.Printf("cache size: %d\n", game.TheCache.Size())

	fmt.Printf("Min Total Turns: %d, %+v\n", minTotalTurns, minTotalTurnsMotivation)
	fmt.Printf("Min Total Pain: %f, %+v\n", minTotalPain, minTotalPainMotivation)
}
