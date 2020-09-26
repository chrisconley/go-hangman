package play

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestCalculateStats(t *testing.T) {
	turns := []int{1, 2, 3, 5, 11, 7}
	stats := CalculateStats(turns)

	test.AssertDeepEqual(stats.GameCount, 6, t)
	test.AssertDeepEqual(stats.TotalTurns, 29, t)
	test.AssertAlmostEqual(stats.AverageTurns, 4.8333, 4, t)
	test.AssertDeepEqual(stats.MaxTurns, 11, t)
}
