package game

import (
	"encoding/binary"
	"fmt"

	"github.com/Workiva/go-datastructures/bitarray"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/mitchellh/hashstructure"
)

type BPossibleWorlds struct {
	Bits bitarray.BitArray
}

func (pw *BPossibleWorlds) GetHash() []byte {
	nums := []byte{}
	iter := pw.Bits.Blocks()
	for iter.Next() {
		_, block := iter.Value()
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(block))
		nums = append(nums, b...)
	}
	// fmt.Printf("nums length: %d\n", len(nums))
	return nums
}

func (pw *BPossibleWorlds) GetBit(index uint64) bool {
	bool, err := pw.Bits.GetBit(index)
	if err != nil {
		panic(err)
	}
	return bool
}

func (pw *BPossibleWorlds) SetBit(index uint64) {
	err := pw.Bits.SetBit(index)
	if err != nil {
		panic(err)
	}
}

func (pw *BPossibleWorlds) GetName(index uint64) string {
	return BUniverse[index].Name
}

func (pw *BPossibleWorlds) ToNums() []uint64 {
	nums := pw.Bits.ToNums()
	return nums[:len(nums)-1]
}

func NewPossibleWorlds(length int) BPossibleWorlds {
	length64 := uint64(length)
	pw := BPossibleWorlds{}
	pw.Bits = bitarray.NewBitArray(length64 + 1)
	pw.Bits.Reset()
	pw.Bits.SetBit(length64)
	return pw
}

func BuildActionWorldHash(action Action, worlds BPossibleWorlds) string {
	actionHash, err := hashstructure.Hash(action, nil)
	if err != nil {
		panic(err)
	}
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, actionHash)
	h := worlds.GetHash()
	h = append(h, b...)
	// hash, err := hashstructure.Hash(sha256.Sum256(h), nil)
	// if err != nil {
	// 	panic(err)
	// }
	return string(h)

}

func LogCacheWorlds(prefix string, hash string, action Action, response string, bpw BPossibleWorlds) {
	fmt.Printf(
		// "UPDATING PERSP| hash:%d, action:%s, response:%s, world:%v\n",
		"%s | hash:%s, action:%s, response:%s, worlds:%+v ",
		prefix,
		hash,
		action,
		response,
		bpw,
	)
	nums := bpw.ToNums()
	if len(nums) < 100 {
		fmt.Printf(
			// "UPDATING PERSP| hash:%d, action:%s, response:%s, world:%v\n",
			"%+v\n",
			bpw.ToNums(),
		)
	} else {
		fmt.Printf("\n")

	}

}

type BWorld struct {
	Index uint64
	Name  string
}

var BUniverse = []BWorld{}
var BPossibleWorldsStore = make(map[string]BPossibleWorlds)
var BBPossibleWorldsStore = make(map[string]map[string]BPossibleWorlds)
var BCountsStore = make(map[string]utils.GroupedIncrementer)
var BRewardsStore = make(map[string]utils.GroupedIncrementer)

func InitUniverse(lines []string) BPossibleWorlds {
	universeLength := len(lines)
	BUniverse = make([]BWorld, universeLength)
	pw := NewPossibleWorlds(universeLength)
	for index, name := range lines {
		BUniverse[index] = BWorld{Index: uint64(index), Name: name}
		pw.SetBit(uint64(index))
	}
	return pw
}

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
