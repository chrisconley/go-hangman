package game

import (
	"github.com/chrisconley/go-hangman/internal/mth"
)

type Predictions map[Action]Prediction

type Prediction struct {
	Action  Action
	Reward  float64
	Entropy float64
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

func BBuildPredictions(actions []Action, bpw BPossibleWorlds) Predictions {
	predictions := Predictions{}
	for _, action := range actions {
		hash := BuildActionWorldHash(action, bpw)
		prediction := Prediction{Action: action}
		counts := BCountsStore[hash]
		rewards := BRewardsStore[hash]
		prediction.Entropy = mth.CalcEntropy(
			counts.Total(),
			counts.All(),
		)
		prediction.Reward = rewards.Total() / counts.Total()
		predictions.AddPrediction(action, prediction)
	}
	return predictions
}
