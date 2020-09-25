package game

import (
	"github.com/chrisconley/go-hangman/internal/mth"
	"github.com/chrisconley/go-hangman/internal/utils"
)

func ExploreReactions(possibleWorlds BPossibleWorlds, actions []Action, act Act) {
	worldsHash := possibleWorlds.GetHash()
	worldIndices := possibleWorlds.Indices()
	for _, action := range actions {
		hash := BuildActionWorldHash(action, worldsHash)
		if _, ok := BBPossibleWorldsStore[hash]; !ok {
			counts, rewards, responseWorlds := ExploreReactionsForAction(
				worldIndices,
				action,
				act,
			)
			BCountsStore[hash] = counts
			BRewardsStore[hash] = rewards
			BBPossibleWorldsStore[hash] = responseWorlds
		}
	}
}

func ExploreReactionsForAction(
	possibleWorldIndices []uint64,
	action Action,
	act Act,
) (
	utils.GroupedIncrementer,
	utils.GroupedIncrementer,
	map[string]BPossibleWorlds,
) {
	counts := utils.NewGroupedIncrementer()
	rewards := utils.NewGroupedIncrementer()
	responseWorlds := make(map[string]BPossibleWorlds)
	for _, worldIndex := range possibleWorldIndices {
		world := BUniverse.Worlds[worldIndex]
		reaction := act(World(world.Name), action)
		response := string(reaction.Response)
		if _, ok := responseWorlds[response]; !ok {
			responseWorlds[response] = NewPossibleWorlds(len(BUniverse.Worlds))
		}
		counts.Increment(response, 1.0)
		rewards.Increment(response, reaction.Reward)
		worlds := responseWorlds[response]
		worlds.SetBit(world.Index)

	}
	return counts, rewards, responseWorlds
}

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
	worldsHash := bpw.GetHash()
	for _, action := range actions {
		hash := BuildActionWorldHash(action, worldsHash)
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
