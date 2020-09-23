package game

import (
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/Workiva/go-datastructures/bitarray"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/mitchellh/hashstructure"
)

type World string
type Action string
type Response string
type Act func(w World, action Action) (reaction Reaction)
type GetActions func(w PossibleWorlds, aa AvailableActions) []Action

type Agreement struct {
	GetActions GetActions
	Act        Act
}

type Reaction struct {
	Action   Action
	Response Response
	Reward   float64
}

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
	return BUniverse.Worlds[index].Name
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
	return string(h)

}

func LogCacheWorlds(prefix string, hash string, action Action, response string, bpw BPossibleWorlds) {
	fmt.Printf(
		"%s | hash:%s, action:%s, response:%s, worlds:%+v ",
		prefix,
		hash,
		action,
		response,
		bpw,
	)
}

type BWorld struct {
	Index uint64
	Name  string
}

type Universe struct {
	Worlds          []BWorld
	Actions         AvailableActions
	PossibleWorlds  PossibleWorlds
	BPossibleWorlds BPossibleWorlds
}

var BUniverse = Universe{}
var BPossibleWorldsStore = make(map[string]BPossibleWorlds)
var BBPossibleWorldsStore = make(map[string]map[string]BPossibleWorlds)
var BCountsStore = make(map[string]utils.GroupedIncrementer)
var BRewardsStore = make(map[string]utils.GroupedIncrementer)

func InitUniverse(lines []string) BPossibleWorlds {
	universeLength := len(lines)
	BUniverse.Worlds = make([]BWorld, universeLength)
	pw := NewPossibleWorlds(universeLength)
	for index, name := range lines {
		BUniverse.Worlds[index] = BWorld{Index: uint64(index), Name: name}
		pw.SetBit(uint64(index))
	}
	BUniverse.PossibleWorlds = BuildPossibleWorlds(lines)
	BUniverse.Actions = BuildAvailableActions(BUniverse.PossibleWorlds)
	BUniverse.BPossibleWorlds = pw
	return pw
}

type PossibleWorlds struct {
	Worlds []World
	key    string
}

func BuildPossibleWorlds(lines []string) PossibleWorlds {
	sort.SliceStable(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})
	key := ""
	worlds := []World{}
	for _, line := range lines {
		key += string(line)
		worlds = append(worlds, World(line))
	}
	return PossibleWorlds{Worlds: worlds, key: key}
}

func BuildPossibleWorldsDeprecated(worlds []World) PossibleWorlds {
	sort.SliceStable(worlds, func(i, j int) bool {
		return worlds[i] < worlds[j]
	})
	key := ""
	for _, world := range worlds {
		key += string(world)
	}
	return PossibleWorlds{Worlds: worlds, key: key}
}

func (w *PossibleWorlds) Key() string {
	return w.key
}

func (pw *PossibleWorlds) Len() int {
	return len(pw.Worlds)
}

type AvailableActions struct {
	worlds       []World
	actions      []Action
	worldActions []Action
	worldSize    int
}

// TODO: Just move this to universe creation then the performance shouldn't matter
// because it'll get called once instead of every new game with NewPerspective
func BuildAvailableActions(worlds PossibleWorlds) AvailableActions {
	aa := AvailableActions{
		actions:      []Action{},
		worldActions: make([]Action, worlds.Len()),
		worldSize:    len(worlds.Worlds[0]),
	}

	keys := make(map[int32]bool)
	for index, world := range worlds.Worlds {
		worldAction := Action(world)
		aa.worldActions[index] = worldAction
		for _, runeValue := range string(world) {
			if _, ok := keys[runeValue]; !ok {
				keys[runeValue] = true
				aa.actions = append(aa.actions, Action(runeValue))
			}
		}
	}

	return aa
}

func (aa *AvailableActions) AllActions() []Action {
	return append(aa.actions, aa.worldActions...)
}

func (aa *AvailableActions) Actions() []Action {
	return append(aa.actions)
}

func (aa *AvailableActions) WorldActions() []Action {
	return append(aa.worldActions)
}

func (aa *AvailableActions) WorldActionsLength() int {
	return len(aa.worldActions)
}

func (aa *AvailableActions) ActionsLength() int {
	return len(aa.actions)
}
