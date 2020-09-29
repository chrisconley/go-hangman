package game

import (
	"fmt"

	"github.com/chrisconley/go-hangman/internal/mth"
	"github.com/chrisconley/go-hangman/internal/utils"
)

var ExploreCount = 0

func ExploreReactions(possibleWorlds BPossibleWorlds, actions []Action, act Act) {
	ExploreCount++
	worldsHash := possibleWorlds.GetHash()
	worldIndices := possibleWorlds.Indices()
	for _, action := range actions {
		hash := BuildActionWorldHash(action, worldsHash)
		if _, ok := BBPossibleWorldsStore[hash]; !ok {
			actionMatchesAPossibleWorld, counts, rewards, responseWorlds := ExploreReactionsForAction(
				worldIndices,
				action,
				act,
			)
			BCountsStore[hash] = counts
			BRewardsStore[hash] = rewards
			BBPossibleWorldsStore[hash] = responseWorlds
			BActionMatchesStore[hash] = actionMatchesAPossibleWorld
			fmt.Printf("SETTING RESPONSE WORLDS (action: %s, len:%d): ", action, len(responseWorlds))
			// fmt.Printf("%+v\n", responseWorlds)
		} else {
			fmt.Printf("REACTIONS ARE CACHED\n")
		}
	}
}

func ExploreReactionsForAction(
	possibleWorldIndices []uint64,
	action Action,
	act Act,
) (
	bool,
	utils.GroupedIncrementer,
	utils.GroupedIncrementer,
	map[string]BPossibleWorlds,
) {
	counts := utils.NewGroupedIncrementer()
	rewards := utils.NewGroupedIncrementer()
	responseWorlds := make(map[string]BPossibleWorlds)
	actionMatchesPossibleWorld := false
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

		if string(action) == world.Name {
			actionMatchesPossibleWorld = true
		}

	}
	return actionMatchesPossibleWorld, counts, rewards, responseWorlds
}

type Predictions map[Action]Prediction

type Prediction struct {
	Action                      Action
	Reward                      float64
	Entropy                     float64
	InversedMax                 float64
	ActionMatchesAPossibleWorld bool
}

func (p Predictions) GetEntropy(action string) float64 {
	return p[Action(action)].Entropy
}

func (p Predictions) GetReward(action string) float64 {
	return p[Action(action)].Reward
}

func (p Predictions) GetInversedMax(action string) float64 {
	return p[Action(action)].InversedMax
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
		prediction.Reward = rewards.Total() / counts.Total()
		prediction.Entropy = mth.CalcEntropy(counts.Total(), counts.All())
		prediction.InversedMax = mth.FindInversedMax(counts.Total(), counts.All())
		prediction.ActionMatchesAPossibleWorld = BActionMatchesStore[hash]
		predictions.AddPrediction(action, prediction)
	}
	fmt.Printf("PREDICTIONS LENGTH: %+v\n", len(predictions))
	for _, prediction := range predictions {
		fmt.Printf("PREDICTION: %+v\n", prediction)

	}
	return predictions
}
