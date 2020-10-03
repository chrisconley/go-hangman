package play

import "github.com/chrisconley/go-hangman/internal/game"

type Stats struct {
	GameCount    int
	TotalTurns   int
	AverageTurns float64
	MaxTurns     int
}

func CalculateStats(games []game.GamePlay) Stats {
	stats := Stats{}
	stats.TotalTurns = 0
	stats.MaxTurns = 0
	stats.GameCount = len(games)
	for _, gameplay := range games {
		stats.TotalTurns = stats.TotalTurns + gameplay.TurnCount
		if gameplay.TurnCount > stats.MaxTurns {
			stats.MaxTurns = gameplay.TurnCount
		}
	}

	stats.AverageTurns = float64(stats.TotalTurns) / float64(stats.GameCount)
	return stats
}
