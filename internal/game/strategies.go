package game

import (
	"math"
)

type ResponseCounts map[Response]*Counts

type Counts struct {
	Weight float32
	Count  int
}
type Reward struct {
	Utility float64
	Regret  float64
}
type ReactionCounts struct {
	ResponseCounts ResponseCounts
	Reward         Reward
}
type ReactionCountsByAction map[Action]ReactionCounts

func (p ReactionCounts) GetTotal() int {
	total := int(0)
	for _, count := range p.ResponseCounts {
		total += count.Count
	}
	return total
}

func ExploreReactions(possibleWorlds PossibleWorlds, actions PossibleActions, agreement Agreement) ReactionCountsByAction {
	reactionCounts := ReactionCountsByAction{}
	for _, action := range actions {
		reactionCounts[action] = ReactionCounts{ResponseCounts: ResponseCounts{}, Reward: Reward{0, 0}}
		for _, possibleWorld := range possibleWorlds {
			reaction := agreement(possibleWorld.World, action)
			_, ok := reactionCounts[action].ResponseCounts[reaction.Response]
			if ok {
				reactionCounts[action].ResponseCounts[reaction.Response].Count += 1
			} else {
				reactionCounts[action].ResponseCounts[reaction.Response] = &Counts{Weight: 0, Count: 1}
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

func calcEntropy(total int, reactionCounts ReactionCounts) float64 {
	entropy := 0.0
	for _, count := range reactionCounts.ResponseCounts {
		entropy += LogProbability(float64(count.Count) / float64(total))
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
