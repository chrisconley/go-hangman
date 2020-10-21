package main

import (
	"encoding/json"
	"fmt"
	"sync"

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

func TryTry(motivation game.Motivation, lines []string, c chan optimizationResults) {
	// limit := 1000
	limit := len(lines)
	games := play.PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := play.CalculateStats(games)

	results := optimizationResults{Motivation: motivation, GameStats: stats}
	js, err := json.Marshal(results)
	if err != nil {
		fmt.Println(js)
	}
	fmt.Println(string(js))

	c <- results
}

func handleResults(input chan optimizationResults, output chan []optimizationResults, wg *sync.WaitGroup) {
	var results []optimizationResults
	for trialResult := range input {
		results = append(results, trialResult)
		wg.Done()
	}
	output <- results
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

	var wg sync.WaitGroup
	limit := 3
	trialResult := make(chan optimizationResults)
	statsOutput := make(chan []optimizationResults)
	go handleResults(trialResult, statsOutput, &wg)
	for _, trial := range trials[:limit] {
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
		wg.Add(1)
		go TryTry(motivation, lines, trialResult)
	}
	wg.Wait()
	close(trialResult)
	allResults := <-statsOutput

	minTotalTurnsMotivation := game.Motivation{}
	minTotalTurns := 1000000

	minTotalPainMotivation := game.Motivation{}
	minTotalPain := 1000000.0
	for _, result := range allResults {
		stats := result.GameStats
		motivation := result.Motivation
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
