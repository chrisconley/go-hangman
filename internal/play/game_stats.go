package play

import "github.com/chrisconley/go-hangman/internal/game"

type Stats struct {
	GameCount    int
	TotalTurns   int
	AverageTurns float64
	MaxTurns     int
}

func CalculateStats(turns []game.GamePlay) Stats {
	stats := Stats{}
	stats.TotalTurns = 0
	stats.MaxTurns = 0
	stats.GameCount = len(turns)
	for _, reactions := range turns {
		stats.TotalTurns = stats.TotalTurns + reactions.TurnCount
		if reactions.TurnCount > stats.MaxTurns {
			stats.MaxTurns = reactions.TurnCount
		}
	}

	stats.AverageTurns = float64(stats.TotalTurns) / float64(stats.GameCount)
	return stats
}
