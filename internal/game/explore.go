package game

import (
	"github.com/chrisconley/go-hangman/internal/mth"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/segmentio/fasthash/fnv1a"
)

func ExploreReactions(
	gameplay GamePlay,
	possibleWorlds PossibleWorlds,
	actions []Action,
	act Act,
) Predictions {
	// If we've already explored the first action,
	// that means we've explored all actions already and can skip.
	worldsHash := possibleWorlds.GetHash()
	hash := BuildActionWorldHash(actions[0], worldsHash)
	predictions := Predictions{PredictionsMap: PredictionsMap{lookup: make(map[uint64]uint32)}}
	// if ok := TheCache.Exists(hash); !ok {
	if _, ok := PredictionCache[hash]; !ok {
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
			prediction := Prediction{Action: action}
			prediction.Reward = rewards.Total() / counts.Total()
			prediction.Entropy = mth.CalcEntropy(counts.Total(), counts.All())
			prediction.PartsCount = float64(len(counts.All()))
			prediction.InversedMax = mth.FindInversedMax(counts.Total(), counts.All())
			prediction.ActionMatchesAPossibleWorld = actionMatchesAPossibleWorld
			predictions.AddPrediction(string(action), prediction)
			PredictionCache[hash] = prediction
			TheCache.Write(
				hash,
				responseWorlds,
				counts,
				rewards,
				actionMatchesAPossibleWorld,
			)
		}
	} else {
		for _, action := range actions {
			hash := BuildActionWorldHash(action, worldsHash)
			predictions.AddPrediction(string(action), PredictionCache[hash])
		}

	}
	return predictions
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

type Predictions struct {
	PredictionsMap PredictionsMap
}

type Prediction struct {
	Action                      Action
	Reward                      float64
	Entropy                     float64
	InversedMax                 float64
	PartsCount                  float64
	ActionMatchesAPossibleWorld bool
}

func (p Predictions) Get(action string) Prediction {
	hash := fnv1a.Init64
	hash = fnv1a.AddString64(hash, action)
	return p.PredictionsMap.Get(hash)
}

func (p Predictions) GetEntropy(action string) float64 {
	hash := fnv1a.Init64
	hash = fnv1a.AddString64(hash, action)
	return p.PredictionsMap.Get(hash).Entropy
}

func (p Predictions) GetReward(action string) float64 {
	hash := fnv1a.Init64
	hash = fnv1a.AddString64(hash, action)
	return p.PredictionsMap.Get(hash).Reward
}

func (p Predictions) GetInversedMax(action string) float64 {
	hash := fnv1a.Init64
	hash = fnv1a.AddString64(hash, action)
	return p.PredictionsMap.Get(hash).InversedMax
}

func (p Predictions) GetPartsCount(action string) float64 {
	hash := fnv1a.Init64
	hash = fnv1a.AddString64(hash, action)
	return p.PredictionsMap.Get(hash).PartsCount
}

func (p *Predictions) AddPrediction(action string, prediction Prediction) {
	hash := fnv1a.Init64
	hash = fnv1a.AddString64(hash, action)
	p.PredictionsMap.Set(hash, prediction)
}

func (p *Predictions) Size() int {
	return int(p.PredictionsMap.Size())
}

func NewPredictions() Predictions {
	p := Predictions{}
	p.PredictionsMap = PredictionsMap{lookup: make(map[uint64]uint32)}
	return p
}

type PredictionsMap struct {
	Values []Prediction
	size   uint32
	lookup map[uint64]uint32
}

func (m *PredictionsMap) Set(group uint64, value Prediction) {
	index := m.findOrInitializeGroup(group)
	m.Values[index] = value
}

func (m *PredictionsMap) Get(group uint64) Prediction {
	index := m.findOrInitializeGroup(group)
	return m.Values[index]
}

func (m *PredictionsMap) Exists(group uint64) bool {
	_, ok := m.lookup[group]
	return ok
}

func (m *PredictionsMap) Size() uint32 {
	return m.size
}

func (m *PredictionsMap) findOrInitializeGroup(group uint64) uint32 {
	index, ok := m.lookup[group]
	if !ok {
		m.Values = append(m.Values, Prediction{})
		m.size++
		m.lookup[group] = m.size - 1
		return m.size - 1
	}
	return index
}

var PredictionsCount uint64
var PredictionCache = make(map[uint64]Prediction)

func BBuildPredictions(actions []Action, bpw PossibleWorlds) Predictions {
	worldsHash := bpw.GetHash()
	predictions := Predictions{PredictionsMap: PredictionsMap{lookup: make(map[uint64]uint32)}}
	hash := BuildActionWorldHash(actions[0], worldsHash)
	if _, ok := PredictionCache[hash]; !ok {
		for _, action := range actions {
			hash := BuildActionWorldHash(action, worldsHash)
			PredictionsCount++
			prediction := Prediction{Action: action}
			counts := TheCache.GetCounts(hash)
			rewards := TheCache.GetRewards(hash)
			prediction.Reward = rewards.Total() / counts.Total()
			prediction.Entropy = mth.CalcEntropy(counts.Total(), counts.All())
			prediction.PartsCount = float64(len(counts.All()))
			prediction.InversedMax = mth.FindInversedMax(counts.Total(), counts.All())
			prediction.ActionMatchesAPossibleWorld = TheCache.GetActionMatches(hash)
			predictions.AddPrediction(string(action), prediction)
		}
	} else {
		for _, action := range actions {
			hash := BuildActionWorldHash(action, worldsHash)
			predictions.AddPrediction(string(action), PredictionCache[hash])
		}
	}
	return predictions

}
