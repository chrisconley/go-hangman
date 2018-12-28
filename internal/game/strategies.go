package game

import (
	"math"
)

type Count int

type ReactionCounts map[Reaction]Count
type ReactionCountsByAction map[Action]ReactionCounts

func (p ReactionCounts) GetTotal() Count {
	total := Count(0)
	for _, Count := range p {
		total += Count
	}
	return total
}

func CountReactions(possibleWorlds PossibleWorlds, actions PossibleActions, agreement Agreement) ReactionCountsByAction {
	reactionCounts := ReactionCountsByAction{}
	for _, action := range actions {
		reactionCounts[action] = ReactionCounts{}
		for _, possibleWorld := range possibleWorlds {
			reaction := agreement(possibleWorld.World, action)
			_, ok := reactionCounts[action][reaction]
			if ok {
				reactionCounts[action][reaction]++
			} else {
				reactionCounts[action][reaction] = 1
			}
		}
	}
	return reactionCounts
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

func calcEntropy(total Count, reactionCounts ReactionCounts) float64 {
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
