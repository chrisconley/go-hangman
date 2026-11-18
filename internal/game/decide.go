package game

import "math"

func Decide(predictions Predictions, motivation Motivation) Prediction {
	maxAction := DecideAction(predictions.PredictionsMap.Values, motivation)
	return predictions.Get(maxAction.String())
}

// DecideAction picks the guess Decide picks, from a plain list of scores.
func DecideAction(predictions []Prediction, motivation Motivation) Action {
	maxWeight := -1.0
	maxAction := NewAction("")
	maxIsValid := false
	if len(predictions) == 1 {
		panic("There is only one prediction for the decision.")
	}
	if len(predictions) == 0 {
		panic("There are no predictions for the decision.")
	}
	for _, prediction := range predictions {
		weight := math.Pow(prediction.Reward, motivation.Reward) *
			math.Pow(prediction.Entropy, motivation.Entropy) *
			math.Pow(prediction.InversedMax, motivation.MiniMax) *
			math.Pow(prediction.PartsCount, motivation.MostParts)

		isValid := prediction.ActionMatchesAPossibleWorld

		// Decide whether or not to swap out maxAction with the new action
		if weight == maxWeight && isValid == maxIsValid && prediction.Action.String() < maxAction.String() {
			// 	// If the new weight is equal to current max
			// 	// AND the new action is lexicographically before maxAction, use it
			maxAction = prediction.Action
		} else if weight == maxWeight && isValid && maxIsValid == false {
			maxAction = prediction.Action
			maxIsValid = isValid
		} else if weight > maxWeight {
			// If there is a clear winner, use it
			maxWeight = weight
			maxAction = prediction.Action
			maxIsValid = isValid
		}
	}
	return maxAction
}
