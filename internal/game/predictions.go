package game

import "github.com/chrisconley/go-hangman/internal/mth"

// "fmt"

type Predictions map[Action]Prediction

type Prediction struct {
	Action  Action
	Reward  float64
	Entropy float64
	//time   float32
}

func (p Predictions) GetEntropy(action string) float64 {
	return p[Action(action)].Entropy
}

func (p Predictions) AddPrediction(action Action, prediction Prediction) {
	p[action] = prediction
}

func BuildPredictions(reactionCounter ReactionCounter) Predictions {
	predictions := Predictions{}
	for _, action := range reactionCounter.Actions() {
		// TODO: Refactor this so there is no direct dict/map access on predictions
		prediction := Prediction{Action: action}
		prediction.Entropy = mth.CalcEntropy(reactionCounter.TotalCountForAction(action), reactionCounter.AllCountsForAction(action))
		predictions.AddPrediction(action, prediction)
	}
	return predictions
}
