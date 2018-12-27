package game

import (
	"fmt"
	"math"
)

type count int
type response string
type guess string

type PossibleResponses map[response]count
type PotentialOutcomes map[guess]PossibleResponses
type GuessWeights map[guess]float64

// LogProbability computes -x * log2(x)
func LogProbability(probability float64) float64 {
	if probability == 0 {
		return 0.0
	}
	return -probability * (math.Log2(probability))
}

func Entropy(potentialOutcomes PotentialOutcomes) GuessWeights {
	guessWeights := GuessWeights{}
	// guessWeights["a"] = 1.0
	for guess, possibleResponses := range potentialOutcomes {
		guessWeights[guess] = singleEntropy(possibleResponses)
	}
	for guess, weight := range guessWeights {
		fmt.Printf("\nguess %s | weight %f", guess, weight)
	}
	fmt.Printf("\nhi")
	return guessWeights
}

func singleEntropy(possibleResponses PossibleResponses) float64 {
	total := count(0)
	for _, count := range possibleResponses {
		total += count
	}
	fmt.Printf("\ntotal $f", total)
	weight := 0.0
	for _, count := range possibleResponses {
		weight += LogProbability(float64(count) / float64(total))
		fmt.Printf("\nweight $f", weight)
	}
	return weight
}
