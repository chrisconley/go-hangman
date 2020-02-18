package mth

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestLogProbability(t *testing.T) {
	test.AssertAlmostEqual(LogProbability(0.0), 0.0, 4, t)
	test.AssertAlmostEqual(LogProbability(0.5), 0.5, 4, t)
	test.AssertAlmostEqual(LogProbability(1.0), 0.0, 4, t)
	test.AssertAlmostEqual(LogProbability(0.01), 0.0664, 4, t)
	test.AssertAlmostEqual(LogProbability(0.99), 0.0144, 4, t)
}

func TestEntropy(t *testing.T) {
	counts := []int{4, 1, 1, 2}
	actual := CalcEntropy(8, counts)
	test.AssertAlmostEqual(actual, 1.75, 4, t)
}
