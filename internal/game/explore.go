package game

import (
	"github.com/chrisconley/go-hangman/internal/mth"
	"github.com/chrisconley/go-hangman/internal/utils"
)

func ExploreReactions(gameplay GamePlay, possibleWorlds PossibleWorlds, actions []Action, act Act) {
	// If we've already explored the first action,
	// that means we've explored all actions already and can skip.
	worldsHash := possibleWorlds.GetHash()
	hash := BuildActionWorldHash(actions[0], worldsHash)
	if ok := TheCache.Exists(hash); !ok {
		// Only get Indices once to improve performance
		worldIndices := possibleWorlds.Indices()
		for _, action := range actions {
			hash := BuildActionWorldHash(action, worldsHash)
			actionMatchesAPossibleWorld, counts, rewards, responseWorlds := ExploreReactionsForAction(
				gameplay,
				worldIndices,
				action,
				act,
			)
			TheCache.Write(
				hash,
				responseWorlds,
				counts,
				rewards,
				actionMatchesAPossibleWorld,
			)
		}
	}
}

func ExploreReactionsForAction(
	gameplay GamePlay,
	possibleWorldIndices []uint64,
	action Action,
	act Act,
) (
	bool,
	utils.GroupedIncrementer,
	utils.GroupedIncrementer,
	map[string]PossibleWorlds,
) {
	counts := utils.NewGroupedIncrementer()
	rewards := utils.NewGroupedIncrementer()
	responseWorlds := make(map[string]PossibleWorlds)
	actionMatchesPossibleWorld := false
	for _, worldIndex := range possibleWorldIndices {
		world := BUniverse.GetWorldByIndex(worldIndex)
		reaction := act(gameplay, World(world.Name), action)
		response := string(reaction.Response)
		if _, ok := responseWorlds[response]; !ok {
			responseWorlds[response] = NewPossibleWorlds(BUniverse.Size)
		}
		counts.Increment(response, 1.0)
		rewards.Increment(response, reaction.Reward)
		worlds := responseWorlds[response]
		worlds.Add(world.Index)
		responseWorlds[response] = worlds

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

func BBuildPredictions(actions []Action, bpw PossibleWorlds) Predictions {
	predictions := Predictions{}
	worldsHash := bpw.GetHash()
	for _, action := range actions {
		hash := BuildActionWorldHash(action, worldsHash)
		prediction := Prediction{Action: action}
		counts := TheCache.GetCounts(hash)
		rewards := TheCache.GetRewards(hash)
		prediction.Reward = rewards.Total() / counts.Total()
		prediction.Entropy = mth.CalcEntropy(counts.Total(), counts.All())
		prediction.InversedMax = mth.FindInversedMax(counts.Total(), counts.All())
		prediction.ActionMatchesAPossibleWorld = TheCache.GetActionMatches(hash)
		predictions.AddPrediction(action, prediction)
	}
	return predictions
}
