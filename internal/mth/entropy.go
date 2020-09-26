package mth

import (
	"math"
)

// FindInversedMax finds the max count then substracts it from the total
// to get the "inversed" maximum.
// We do this because we want to minimize the max for the
// MiniMax motivation (in Knuth's paper), but our decision algorithm
// finds the choices that maximize motivation. So in the end, instead of
// minimizing the max, we can maximize the inversed max to stay consistent
// with the other motivations.
func FindInversedMax(total float64, counts []float64) float64 {
	maxCount := 0.0
	for _, count := range counts {
		if count > maxCount {
			maxCount = count
		}
	}
	return total - maxCount
}

// CalcEntropy calculates the information entropy contained in the
// slice of counts.
func CalcEntropy(total float64, counts []float64) float64 {
	entropy := 0.0
	for _, count := range counts {
		entropy += logProbability(count / total)
	}
	return entropy
}

// logProbability computes -x * log2(x)
func logProbability(probability float64) float64 {
	if probability == 0 {
		return 0.0
	}
	return -probability * (math.Log2(probability))
}
