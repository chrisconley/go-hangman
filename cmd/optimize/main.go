package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/coocood/freecache"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/chrisconley/go-hangman/internal/wordle"
	"github.com/pkg/profile"
)

type optimizationResults struct {
	Motivation game.Motivation
	GameStats  play.Stats
}

func TryTry(motivation game.Motivation, agreement game.Agreement, lines []string, c chan optimizationResults, wg *sync.WaitGroup) {
	limit := len(lines)
	games := play.PlayGames(agreement, motivation, lines, limit)
	stats := play.CalculateStats(games, true)

	results := optimizationResults{Motivation: motivation, GameStats: stats}
	js, err := json.Marshal(results)
	if err != nil {
		fmt.Println(js)
	}
	fmt.Println(string(js))

	wg.Done()
	c <- results
}

// parseWeights reads the -weights list into trials, each in the order BuildTrials uses: entropy,
// minimax, mostparts, reward.
func parseWeights(list string) [][]float64 {
	trials := [][]float64{}
	for _, weighting := range strings.Split(list, ",") {
		parts := strings.Split(weighting, ":")
		if len(parts) != 4 {
			panic(fmt.Sprintf("-weights: %q should be reward:entropy:minimax:mostparts", weighting))
		}
		values := make([]float64, 4)
		for i, part := range parts {
			value, err := strconv.ParseFloat(part, 64)
			if err != nil {
				panic(fmt.Sprintf("-weights: %q is not a number", part))
			}
			values[i] = value
		}
		trials = append(trials, []float64{values[1], values[2], values[3], values[0]})
	}
	return trials
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
	lines, err := utils.ReadLinesFromStdin()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	limit := flag.Int("limit", -1, "Number of trials to run. Defaults to all generated")
	gameString := flag.String("game", "none", "What game to play. Valid games: hangman, mastermind, wordle. Required")
	wordleGuesses := flag.String("wordle.guesses", "", "With -game wordle: file of allowed guesses, one per line. The answers come on stdin.")
	trialStepSize := flag.Float64("trial.stepSize", 0.5, "What step size to take when building brute force trials")
	weights := flag.String("weights", "", "Play only these weightings instead of every combination at -trial.stepSize: a comma-separated list, each written reward:entropy:minimax:mostparts, e.g. 0:1:0:0,0:0.7:0.15:0.15")
	profiling := flag.String("profiling", "off", "Set 'mem' or 'cpu' profiling. Defaults to off.")
	profilingProfilePath := flag.String("profiling.profilePath", "/tmp/profiling", "Path where profile should be written")
	cacheSize := flag.Int("cacheSize", 1000, "Cache size in MB. Defaults to 1000MB (1GB).")
	flag.Parse()

	fmt.Println("Cache Size in MB: ", *cacheSize)
	game.ThePredictionsFreeCache = freecache.NewCache(*cacheSize * 1024 * 1024)
	debug.SetGCPercent(1)

	agreement := hangman.Agreement
	if *gameString == "hangman" {
		agreement = hangman.Agreement
	} else if *gameString == "mastermind" {
		agreement = mastermind.Agreement
	} else if *gameString == "wordle" {
		agreement = wordle.NewAgreement(utils.ReadLinesFromFilePath(*wordleGuesses))
	} else {
		panic("A game must be specified.")
	}

	if *profiling == "cpu" {
		defer profile.Start(
			profile.ProfilePath(*profilingProfilePath)).Stop()
	} else if *profiling == "mem" {
		defer profile.Start(
			// profile.MemProfile,
			profile.MemProfileHeap,
			profile.MemProfileRate(4096*10),
			profile.ProfilePath(*profilingProfilePath),
		).Stop()
	}

	game.BUniverse = game.InitUniverse(lines)
	var trials [][]float64
	if *weights != "" {
		trials = parseWeights(*weights)
	} else {
		trials = play.BuildTrials(0.0, 1.0, *trialStepSize, 4)
	}
	trialLimit := 0
	if *limit == -1 {
		trialLimit = len(trials)
	} else {
		trialLimit = *limit
	}
	fmt.Println("Number of trials generated: ", len(trials))
	fmt.Println("Trials to be conducted: ", trialLimit)

	var wg sync.WaitGroup
	trialResult := make(chan optimizationResults)
	statsOutput := make(chan []optimizationResults)
	maxGoroutines := runtime.NumCPU()
	guard := make(chan struct{}, maxGoroutines)
	go handleResults(trialResult, statsOutput, &wg)

	for _, trial := range trials[:trialLimit] {
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
		wg.Add(2)
		guard <- struct{}{}
		go func() {
			TryTry(motivation, agreement, lines, trialResult, &wg)
			<-guard
		}()
	}
	wg.Wait()
	close(trialResult)
	allResults := <-statsOutput

	minTotalTurnsMotivation := game.Motivation{}
	minTotalTurns := 1000000

	minTotalPainMotivation := game.Motivation{}
	minTotalPain := 1000000.0

	minTotalPainRewardRatioMotivation := game.Motivation{}
	minTotalPainRewardRatio := 1000000.0
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
		if stats.TotalPain/stats.TotalReward < minTotalPainRewardRatio {
			minTotalPainRewardRatio = stats.TotalPain / stats.TotalReward
			minTotalPainRewardRatioMotivation = motivation
		}

	}

	fmt.Printf("Min Total Turns: %d, %+v\n", minTotalTurns, minTotalTurnsMotivation)
	fmt.Printf("Min Total Pain: %f, %+v\n", minTotalPain, minTotalPainMotivation)
	fmt.Printf("Min Total Pain/Reward Ratio: %f, %+v\n", minTotalPainRewardRatio, minTotalPainRewardRatioMotivation)

}
