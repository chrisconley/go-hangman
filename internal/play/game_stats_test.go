package play

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/test"
)

func TestCalculateStats(t *testing.T) {
	turns := []game.GamePlay{
		game.GamePlay{TurnCount: 1},
		game.GamePlay{TurnCount: 2},
		game.GamePlay{TurnCount: 3},
		game.GamePlay{TurnCount: 5},
		game.GamePlay{TurnCount: 11},
		game.GamePlay{TurnCount: 7},
	}
	stats := CalculateStats(turns)

	test.AssertDeepEqual(stats.GameCount, 6, t)
	test.AssertDeepEqual(stats.TotalTurns, 29, t)
	test.AssertAlmostEqual(stats.AverageTurns, 4.8333, 4, t)
	test.AssertDeepEqual(stats.MaxTurns, 11, t)
}
