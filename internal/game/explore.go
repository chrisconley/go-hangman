package game

import (
	"bytes"
	"encoding/binary"

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
	predictions := NewPredictions()
	// Only get Indices once to improve performance
	worldIndices := possibleWorlds.Indices()
	for _, action := range actions {
		actionMatchesAPossibleWorld, counts, rewards, _ := ExploreReactionsForAction(
			gameplay,
			worldIndices,
			action,
			act,
		)

		reward, sizes := GroupsInFixedOrder(counts.All(), rewards.All())
		prediction := Prediction{Action: action}
		prediction.Reward = reward / counts.Total()
		prediction.Entropy = mth.CalcEntropy(counts.Total(), sizes)
		prediction.PartsCount = float64(len(counts.All()))
		prediction.InversedMax = mth.FindInversedMax(counts.Total(), counts.All())
		prediction.ActionMatchesAPossibleWorld = actionMatchesAPossibleWorld
		predictions.AddPrediction(action.String(), prediction)

	}
	return predictions

}

// GroupsInFixedOrder returns a guess's total reward and its group sizes, with the groups taken in a
// fixed order (smallest reward first, then smallest group) instead of the order the answers happen to
// come in. Adding the same numbers in a different order can change the last decimal place of the sum, so
// without this two guesses that split the answers into exactly the same groups could score a hair apart,
// and the tie between them would be settled by rounding instead of by Decide's rule.
func GroupsInFixedOrder(sizes []float64, rewards []float64) (float64, []float64) {
	orderedRewards := make([]float64, len(sizes))
	ordered := make([]float64, len(sizes))
	for i := range sizes {
		j := i
		for j > 0 && (rewards[i] < orderedRewards[j-1] || (rewards[i] == orderedRewards[j-1] && sizes[i] < ordered[j-1])) {
			orderedRewards[j], ordered[j] = orderedRewards[j-1], ordered[j-1]
			j--
		}
		orderedRewards[j], ordered[j] = rewards[i], sizes[i]
	}
	reward := 0.0
	for _, r := range orderedRewards {
		reward += r
	}
	return reward, ordered
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
	map[uint64]PossibleWorlds,
) {
	counts := utils.NewGroupedIncrementer()
	rewards := utils.NewGroupedIncrementer()
	responseWorlds := make(map[uint64]PossibleWorlds)
	actionMatchesPossibleWorld := false
	for _, worldIndex := range possibleWorldIndices {
		world := BUniverse.GetWorldByIndex(worldIndex)
		reaction := act(gameplay, world, action)
		responseString := reaction.Response.String()
		responseHash := fnv1a.HashString64(responseString)
		if _, ok := responseWorlds[responseHash]; !ok {
			responseWorlds[responseHash] = NewPossibleWorlds(BUniverse.Size)
		}
		counts.Increment(responseString, 1.0)
		rewards.Increment(responseString, reaction.Reward)
		worlds := responseWorlds[responseHash]
		worlds.Add(worldIndex)
		responseWorlds[responseHash] = worlds

		if action.Matches(world) {
			actionMatchesPossibleWorld = true
		}

	}
	return actionMatchesPossibleWorld, counts, rewards, responseWorlds
}

func GetPossibleResponses(
	gameplay GamePlay,
	possibleWorlds PossibleWorlds,
	actions []Action,
	act Act,
) []string {
	responseSet := map[string]struct{}{}
	responses := []string{}
	worldIndices := possibleWorlds.Indices()
	for _, action := range actions {
		for _, worldIndex := range worldIndices {
			world := BUniverse.GetWorldByIndex(worldIndex)
			reaction := act(gameplay, world, action)
			responseString := reaction.Response.String()
			if _, ok := responseSet[responseString]; !ok {
				responseSet[responseString] = struct{}{}
				responses = append(responses, responseString)
			}
		}
	}
	return responses
}

type Prediction struct {
	Action                      Action
	Reward                      float64
	Entropy                     float64
	InversedMax                 float64
	PartsCount                  float64
	ActionMatchesAPossibleWorld bool
}

func (p Prediction) ToBytes() []byte {
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.LittleEndian, p.Reward)
	binary.Write(&buffer, binary.LittleEndian, p.Entropy)
	binary.Write(&buffer, binary.LittleEndian, p.InversedMax)
	binary.Write(&buffer, binary.LittleEndian, p.PartsCount)
	binary.Write(&buffer, binary.LittleEndian, p.ActionMatchesAPossibleWorld)

	buffer.Write(p.Action.Bytes())
	return buffer.Bytes()
}

func NewPredictionFromBytes(data []byte) Prediction {
	// We make a copy of data to ensure the slice being passed in
	// can be garbage collected.
	predictionBytes := make([]byte, len(data), len(data))
	copy(predictionBytes, data)
	data = nil // make sure we don't accidentally reference data below

	buffer := bytes.NewBuffer(predictionBytes[:33])
	var reward float64
	var entropy float64
	var inversedMax float64
	var partsCount float64
	var actionMatches bool
	binary.Read(buffer, binary.LittleEndian, &reward)
	binary.Read(buffer, binary.LittleEndian, &entropy)
	binary.Read(buffer, binary.LittleEndian, &inversedMax)
	binary.Read(buffer, binary.LittleEndian, &partsCount)
	binary.Read(buffer, binary.LittleEndian, &actionMatches)

	action := NewActionFromBytes(predictionBytes[33:])

	return Prediction{
		Action:                      action,
		Reward:                      reward,
		Entropy:                     entropy,
		InversedMax:                 inversedMax,
		PartsCount:                  partsCount,
		ActionMatchesAPossibleWorld: actionMatches,
	}

}

type Predictions struct {
	PredictionsMap PredictionsMap
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

func (p *Predictions) ToBytes() []byte {
	buffer := bytes.Buffer{}
	var dataLength uint32

	binary.Write(&buffer, binary.LittleEndian, uint32(p.Size()))
	for _, prediction := range p.PredictionsMap.Values {
		data := prediction.ToBytes()
		dataLength = uint32(len(data))
		binary.Write(&buffer, binary.LittleEndian, dataLength)
		buffer.Write(prediction.ToBytes())
	}
	return buffer.Bytes()
}

func NewPredictions() Predictions {
	p := Predictions{}
	p.PredictionsMap = PredictionsMap{lookup: make(map[uint64]uint32)}

	return p
}

func NewPredictionsFromBytes(data []byte) Predictions {
	// We make a copy of data to ensure the slice being passed in
	// can be garbage collected.
	predictionsBytes := make([]byte, len(data), len(data))
	copy(predictionsBytes, data)
	data = nil // make sure we don't accidentally reference data below

	p := NewPredictions()
	pointer := 0
	int32Size := 4
	var size uint32
	var predictionDataLength uint32

	size = binary.LittleEndian.Uint32(predictionsBytes[pointer : pointer+int32Size])
	pointer += int32Size

	for i := 0; i < int(size); i++ {
		predictionDataLength = binary.LittleEndian.Uint32(predictionsBytes[pointer : pointer+int32Size])
		pointer += int32Size

		predictionData := predictionsBytes[pointer : pointer+int(predictionDataLength)]
		pointer += int(predictionDataLength)

		prediction := NewPredictionFromBytes(predictionData)
		p.AddPrediction(prediction.Action.String(), prediction)
	}
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
