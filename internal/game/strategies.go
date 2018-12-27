package game

import (
	"math"
)

type count int

type ReactionCounts map[Reaction]count
type ReactionCountsByAction map[Action]ReactionCounts

func (p ReactionCounts) GetTotal() count {
	total := count(0)
	for _, count := range p {
		total += count
	}
	return total
}

func BuildPredictions(reactionCountsByAction ReactionCountsByAction) Predictions {
	predictions := Predictions{}
	for action, reactionCounts := range reactionCountsByAction {
		predictions[action] = buildActionForecasts(action, reactionCounts)
	}
	return predictions
}

func buildActionForecasts(action Action, reactionCounts ReactionCounts) PredictedActionUtility {
	forecast := PredictedActionUtility{Action: action}
	forecast.Entropy = calcEntropy(reactionCounts.GetTotal(), reactionCounts)
	return forecast
}

func calcEntropy(total count, reactionCounts ReactionCounts) float64 {
	entropy := 0.0
	for _, count := range reactionCounts {
		entropy += LogProbability(float64(count) / float64(total))
	}
	return entropy
}

// LogProbability computes -x * log2(x)
func LogProbability(probability float64) float64 {
	if probability == 0 {
		return 0.0
	}
	return -probability * (math.Log2(probability))
}
