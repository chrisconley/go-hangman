package game

import (
	// "fmt"
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

func (r ReactionCountsByAction) findOrInitializeResponse(action Action, response Response) {
	_, ok := r[action]
	if ok {

	} else {
		r[action] = ReactionCounts{ResponseCounts: ResponseCounts{}, Reward: Reward{0, 0}}
	}
	_, ok = r[action].ResponseCounts[response]
	if ok {

	} else {
		r[action].ResponseCounts[response] = &Counts{Weight: 0, Count: 0}
	}
}

func (r ReactionCountsByAction) GetResponseCount(action Action, response Response) int {
	r.findOrInitializeResponse(action, response)
	return r[action].ResponseCounts[response].Count
}

func (r ReactionCountsByAction) IncrementResponseCount(action Action, response Response) {
	r.findOrInitializeResponse(action, response)
	r[action].ResponseCounts[response].Count += 1
}

func (r ReactionCountsByAction) SetResponseCounts(action Action, response Response, weight float32, count int) {
	r.findOrInitializeResponse(action, response)
	r[action].ResponseCounts[response] = &Counts{Weight: weight, Count: count}
}

func ExploreReactions(possibleWorlds PossibleWorlds, actions PossibleActions, agreement Agreement) ReactionCountsByAction {
	reactionCounts := ReactionCountsByAction{}
	for _, action := range actions {
		for _, possibleWorld := range possibleWorlds {
			reaction := agreement(possibleWorld.World, action)
			reactionCounts.IncrementResponseCount(action, reaction.Response)
		}
	}
	// fmt.Printf("\nREACTIONCOUNTS: %+v", reactionCounts)
	return reactionCounts
}

func BuildPredictions(reactionCountsByAction ReactionCountsByAction) Predictions {
	predictions := Predictions{}
	for action, reactionCounts := range reactionCountsByAction {
		// TODO: Refactor this so there is no direct dict/map access on predictions
		predictions.AddPrediction(action, buildActionForecasts(action, reactionCounts))
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
