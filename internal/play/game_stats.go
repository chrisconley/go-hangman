package play

import "github.com/chrisconley/go-hangman/internal/game"

type Stats struct {
	GameCount     int
	TotalTurns    int
	AverageTurns  float64
	MaxTurns      int
	TotalReward   float64
	AverageReward float64
	TotalPain     float64
	AveragePain   float64
	TurnCounts    []int
}

func CalculateStats(games []game.GamePlay, includeTurnCounts bool) Stats {
	stats := Stats{}
	stats.GameCount = len(games)
	for _, gameplay := range games {
		if includeTurnCounts {
			stats.TurnCounts = append(stats.TurnCounts, gameplay.TurnCount)
		}
		stats.TotalTurns = stats.TotalTurns + gameplay.TurnCount
		stats.TotalReward += gameplay.TotalReward
		stats.TotalPain += gameplay.TotalPain
		if gameplay.TurnCount > stats.MaxTurns {
			stats.MaxTurns = gameplay.TurnCount
		}
	}

	stats.AverageTurns = float64(stats.TotalTurns) / float64(stats.GameCount)
	stats.AverageReward = stats.TotalReward / float64(stats.GameCount)
	stats.AveragePain = stats.TotalPain / float64(stats.GameCount)
	return stats
}
