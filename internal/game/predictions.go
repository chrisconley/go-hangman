package game

import (
	"github.com/chrisconley/go-hangman/internal/mth"
	"github.com/mitchellh/hashstructure"
)

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

func BBuildPredictions(actions []Action, bpw BPossibleWorlds) Predictions {
	predictions := Predictions{}
	for _, action := range actions {
		actionHash, err := hashstructure.Hash(action, nil)
		if err != nil {
			panic(err)
		}
		hash, err := hashstructure.Hash([]uint64{actionHash, bpw.GetHash()}, nil)
		if err != nil {
			panic(err)
		}
		prediction := Prediction{Action: action}
		counts := BCountsStore[hash]
		rewards := BRewardsStore[hash]
		prediction.Entropy = mth.CalcEntropy(
			counts.Total(),
			counts.All(),
		)
		prediction.Reward = rewards.Total() / float64(counts.Total())
		predictions.AddPrediction(action, prediction)
	}
	// fmt.Printf("PREDICTIONS %v\n", predictions)
	return predictions
}
