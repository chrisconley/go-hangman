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

func (p Predictions) GetReward(action string) float64 {
	return p[Action(action)].Reward
}

func (p Predictions) AddPrediction(action Action, prediction Prediction) {
	p[action] = prediction
}

func BuildPredictions(reactionExplorer reactionExplorer) Predictions {
	predictions := Predictions{}
	for _, action := range reactionExplorer.Actions() {
		// TODO: Refactor this so there is no direct dict/map access on predictions
		prediction := Prediction{Action: action}
		prediction.Entropy = mth.CalcEntropy(
			reactionExplorer.ActionTotal(action),
			reactionExplorer.ActionCounts(action),
		)
		prediction.Reward = reactionExplorer.ActionRewardTotal(action) / float64(reactionExplorer.ActionTotal(action))
		predictions.AddPrediction(action, prediction)
	}
	return predictions
}
