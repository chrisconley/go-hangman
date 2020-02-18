package mth

import (
	"math"
)

func CalcEntropy(total int, counts []int) float64 {
	entropy := 0.0
	for _, count := range counts {
		entropy += LogProbability(float64(count) / float64(total))
	}
	return entropy
}

// LogProbability computes -x * log2(x)
func LogProbability(probability float64) float64 {
	if probability == 0 {
		return 0.0
	}
	return -probability * (math.Log2(probability))
}
