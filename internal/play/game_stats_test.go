package play

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/test"
)

func TestCalculateStats(t *testing.T) {
	games := []game.GamePlay{
		game.GamePlay{TurnCount: 1},
		game.GamePlay{TurnCount: 2},
		game.GamePlay{TurnCount: 3},
		game.GamePlay{TurnCount: 5},
		game.GamePlay{TurnCount: 11},
		game.GamePlay{TurnCount: 7},
	}
	stats := CalculateStats(games)

	test.AssertDeepEqual(stats.GameCount, 6, t)
	test.AssertDeepEqual(stats.TotalTurns, 29, t)
	test.AssertAlmostEqual(stats.AverageTurns, 4.8333, 4, t)
	test.AssertDeepEqual(stats.MaxTurns, 11, t)
}

func TestCalculateStatsReward(t *testing.T) {
	g1 := game.GamePlay{}
	g1.AddTurn(game.Reaction{Reward: 0.0})
	g1.AddTurn(game.Reaction{Reward: 0.75})

	g2 := game.GamePlay{}
	g2.AddTurn(game.Reaction{Reward: 1.0})

	g3 := game.GamePlay{}
	g3.AddTurn(game.Reaction{Reward: 0.0})

	g4 := game.GamePlay{}
	g4.AddTurn(game.Reaction{Reward: 0.0})
	g4.AddTurn(game.Reaction{Reward: 0.333333})
	g4.AddTurn(game.Reaction{Reward: 1.0})

	games := []game.GamePlay{g1, g2, g3, g4}
	stats := CalculateStats(games)

	test.AssertAlmostEqual(stats.AverageReward, 0.7708, 4, t)
	test.AssertAlmostEqual(stats.TotalReward, 3.0833, 4, t)

	test.AssertAlmostEqual(stats.AveragePain, 0.9792, 4, t)
	test.AssertAlmostEqual(stats.TotalPain, 3.9167, 4, t)
}
