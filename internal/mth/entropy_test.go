package mth

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

// Test the computation of LogProbility through CalcEntropy
// with very simple cases.
func TestLogProbability(t *testing.T) {
	// 0%
	test.AssertAlmostEqual(CalcEntropy(1, []float64{0.0}), 0.0, 4, t)

	// 50%
	test.AssertAlmostEqual(CalcEntropy(2, []float64{1.0}), 0.5, 4, t)

	// 100%
	test.AssertAlmostEqual(CalcEntropy(1, []float64{1.0}), 0.0, 4, t)

	// 1%
	test.AssertAlmostEqual(CalcEntropy(100, []float64{1.0}), 0.0664, 4, t)

	// 99%
	test.AssertAlmostEqual(CalcEntropy(100, []float64{99.0}), 0.0144, 4, t)
}

func TestEntropy(t *testing.T) {
	counts := []float64{4, 1, 1, 2}
	actual := CalcEntropy(8, counts)
	test.AssertAlmostEqual(actual, 1.75, 4, t)
}

func TestFindMaxInverse(t *testing.T) {
	counts := []float64{1, 1, 5, 1}
	test.AssertAlmostEqual(FindInversedMax(8, counts), 3.0, 4, t)

	counts = []float64{4, 1, 1, 2}
	test.AssertAlmostEqual(FindInversedMax(8, counts), 4.0, 4, t)

	counts = []float64{1, 4, 1, 2}
	test.AssertAlmostEqual(FindInversedMax(8, counts), 4.0, 4, t)

	counts = []float64{2, 3, 1, 2}
	test.AssertAlmostEqual(FindInversedMax(8, counts), 5.0, 4, t)

	counts = []float64{2, 2, 2, 2}
	test.AssertAlmostEqual(FindInversedMax(8, counts), 6.0, 4, t)
}
