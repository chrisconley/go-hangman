package main

import (
	"fmt"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/pkg/profile"
)

func TryTry(motivation game.Motivation, lines []string) play.Stats {
	// limit := 1000
	limit := len(lines)
	fmt.Printf("%+v | ", motivation)
	games := play.PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := play.CalculateStats(games)
	// fmt.Printf("Pain: %f, %+v\n", stats.TotalPain, stats)
	// return float64(stats.TotalPain)
	fmt.Printf("%+v\n", stats)
	// return float64(stats.TotalTurns)
	return stats
}

func main() {
	lines := utils.ReadLinesFromFilePath("./fixtures/commercial-mastermind.txt")
	// lines := utils.ReadLinesFromFilePath("./fixtures/nine-letter-words.txt")
	// lines = lines[:10000]
	game.BUniverse = game.InitUniverse(lines)
	trials := play.BuildTrials(0.0, 1.0, 0.5, 4)
	// fmt.Println(len(trials))
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
	fmt.Printf("predictions size: %d", game.PredictionsCount)

	fmt.Printf("Min Total Turns: %d, %+v\n", minTotalTurns, minTotalTurnsMotivation)
	fmt.Printf("Min Total Pain: %f, %+v\n", minTotalPain, minTotalPainMotivation)
}

// func main() {
// 	lines := utils.ReadLinesFromFilePath("./fixtures/commercial-mastermind.txt")
// 	// lines := utils.ReadLinesFromFilePath("./fixtures/nine-letter-words.txt")
// 	// lines = lines[:10000]
// 	game.BUniverse = game.InitUniverse(lines)

// 	fcn := func(x []float64) float64 {
// 		for _, value := range x {
// 			if value < 0.0 {
// 				return 1000000.00
// 			}
// 		}
// 		motivation := game.Motivation{
// 			// Reward:  0.0,
// 			Entropy: x[0],
// 			// MiniMax:   x[1],
// 			MostParts: x[1],
// 		}
// 		return TryTry(motivation, lines)
// 	}

// 	// grad := func(grad, x []float64) {
// 	// 	fd.Gradient(grad, fcn, x, nil)
// 	// }

// 	// hess := func(h *mat.SymDense, x []float64) {
// 	// 	fd.Hessian(h, fcn, x, nil)
// 	// }

// 	problem := optimize.Problem{
// 		Func: fcn,
// 		// Grad: grad,
// 		// Hess: hess,
// 	}

// 	// method := &optimize.Newton{}
// 	// method := &optimize.NelderMead{}
// 	result, err := optimize.Minimize(
// 		problem,
// 		[]float64{0.001, 0.0001},
// 		nil,
// 		nil,
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println("result:", result.X)

// }
