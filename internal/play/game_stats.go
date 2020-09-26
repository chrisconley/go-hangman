package play

type Stats struct {
	GameCount    int
	TotalTurns   int
	AverageTurns float64
	MaxTurns     int
}

func CalculateStats(turnCounts []int) Stats {
	stats := Stats{}
	stats.TotalTurns = 0
	stats.MaxTurns = 0
	stats.GameCount = len(turnCounts)
	for _, turnCount := range turnCounts {
		stats.TotalTurns = stats.TotalTurns + turnCount
		if turnCount > stats.MaxTurns {
			stats.MaxTurns = turnCount
		}
	}

	stats.AverageTurns = float64(stats.TotalTurns) / float64(stats.GameCount)
	return stats
}
