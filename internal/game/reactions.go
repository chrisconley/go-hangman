package game

import (
	"github.com/Workiva/go-datastructures/bitarray"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/mitchellh/hashstructure"
)

type BPossibleWorlds struct {
	Bits bitarray.BitArray
}

func (pw *BPossibleWorlds) GetHash() uint64 {
	hash, err := hashstructure.Hash(pw, nil)
	if err != nil {
		panic(err)
	}
	return hash
}

func (pw *BPossibleWorlds) GetBit(index uint64) bool {
	bool, err := pw.Bits.GetBit(index)
	if err != nil {
		panic(err)
	}
	return bool
}

func (pw *BPossibleWorlds) GetName(index uint64) string {
	return BUniverse[index].Name
}

func NewPossibleWorlds(length uint64) BPossibleWorlds {
	pw := BPossibleWorlds{}
	pw.Bits = bitarray.NewBitArray(length)
	pw.Bits.Reset()
	return pw
}

type BWorld struct {
	Index uint64
	Name  string
}

var BUniverse = []BWorld{}
var BPossibleWorldsStore = make(map[uint64]BPossibleWorlds)
var BBPossibleWorldsStore = make(map[uint64]map[string]BPossibleWorlds)
var BCountsStore = make(map[uint64]utils.GroupedIncrementer)
var BRewardsStore = make(map[uint64]utils.GroupedIncrementer)

type reactionExplorer struct {
	tracker map[Action]actionData
	actions []Action
}

type actionData struct {
	Responses map[Response][]World
	Counts    *utils.GroupedIncrementer
	Rewards   *utils.GroupedIncrementer
}

func NewReactionExplorer() reactionExplorer {
	re := reactionExplorer{}
	re.tracker = make(map[Action]actionData)
	re.actions = make([]Action, 0)
	return re
}

func (re *reactionExplorer) Add(action Action, response Response, reward float64, world World) {
	re.findOrInitializeResponse(action, response)
	counts := re.tracker[action].Counts
	counts.Increment(string(response), 1.0)
	rewards := re.tracker[action].Rewards
	rewards.Increment(string(response), reward)
	re.tracker[action].Responses[response] = append(re.tracker[action].Responses[response], world)
}

func (re *reactionExplorer) Get(action Action, response Response) []World {
	re.findOrInitializeResponse(action, response)
	return re.tracker[action].Responses[response]
}

func (re *reactionExplorer) ActionCounts(action Action) []float64 {
	re.findOrInitializeAction(action)
	return re.tracker[action].Counts.All()
}

func (re *reactionExplorer) ActionTotal(action Action) float64 {
	re.findOrInitializeAction(action)
	return re.tracker[action].Counts.Total()
}

func (re *reactionExplorer) ActionRewardTotal(action Action) float64 {
	re.findOrInitializeAction(action)
	return re.tracker[action].Rewards.Total()
}

func (re *reactionExplorer) Actions() []Action {
	return re.actions
}

func (re *reactionExplorer) findOrInitializeAction(action Action) {
	_, ok := re.tracker[action]
	if !ok {
		re.actions = append(re.actions, action)
		counts := utils.NewGroupedIncrementer()
		rewards := utils.NewGroupedIncrementer()
		re.tracker[action] = actionData{
			Responses: make(map[Response][]World),
			Counts:    &counts,
			Rewards:   &rewards,
		}
	}
}

func (re *reactionExplorer) findOrInitializeResponse(action Action, response Response) {
	re.findOrInitializeAction(action)
	_, ok := re.tracker[action].Responses[response]
	if !ok {
		re.tracker[action].Responses[response] = make([]World, 0)
	}
}
