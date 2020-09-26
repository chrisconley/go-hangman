package play

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/test"
)

func TestCalculateStats(t *testing.T) {
	turns := []game.Reactions{
		game.Reactions{TurnCount: 1},
		game.Reactions{TurnCount: 2},
		game.Reactions{TurnCount: 3},
		game.Reactions{TurnCount: 5},
		game.Reactions{TurnCount: 11},
		game.Reactions{TurnCount: 7},
	}
	stats := CalculateStats(turns)

	test.AssertDeepEqual(stats.GameCount, 6, t)
	test.AssertDeepEqual(stats.TotalTurns, 29, t)
	test.AssertAlmostEqual(stats.AverageTurns, 4.8333, 4, t)
	test.AssertDeepEqual(stats.MaxTurns, 11, t)
}
