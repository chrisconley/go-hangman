package game

import (
	"encoding/binary"

	"github.com/RoaringBitmap/roaring/roaring64"
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

type Reactions struct {
	Turns     []Reaction
	TurnCount int
}

type PossibleWorlds struct {
	Bits     *roaring64.Bitmap
	count    int
	Capacity int
}

func (pw *PossibleWorlds) GetHash() []byte {
	bytes, _ := pw.Bits.ToBytes()
	return bytes
}

func (pw *PossibleWorlds) Add(index uint64) {
	pw.count += 1
	pw.Bits.Add(index)
}

func (pw *PossibleWorlds) Len() int {
	return pw.count
}

func (pw *PossibleWorlds) Indices() []uint64 {
	nums := pw.Bits.ToArray()
	return nums[:len(nums)-1]
}

func NewPossibleWorlds(length int) PossibleWorlds {
	length64 := uint64(length)
	pw := PossibleWorlds{}
	pw.Bits = roaring64.BitmapOf(length64 + 1)
	pw.count = 0
	pw.Capacity = length
	return pw
}

func BuildActionWorldHash(action Action, worldsHash []byte) string {
	actionHash, err := hashstructure.Hash(action, nil)
	if err != nil {
		panic(err)
	}
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, actionHash)
	worldsHash = append(worldsHash, b...)
	return string(worldsHash)

}

type BWorld struct {
	Index uint64
	Name  string
}

type Universe struct {
	Worlds         []BWorld
	WorldLength    int
	Actions        AvailableActions
	PossibleWorlds PossibleWorlds
	Size           int
}

var BUniverse = Universe{}

type Cache struct{}

var TheCache = Cache{}
var worldsCache = make(map[string]map[string]PossibleWorlds)
var countsCache = make(map[string]utils.GroupedIncrementer)
var rewardsCache = make(map[string]utils.GroupedIncrementer)
var actionMatchesCache = make(map[string]bool)

func (c *Cache) GetNextWorlds(hash string, response string) PossibleWorlds {
	return worldsCache[hash][response]
}

func (c *Cache) Write(
	hash string,
	responseWorlds map[string]PossibleWorlds,
	counts utils.GroupedIncrementer,
	rewards utils.GroupedIncrementer,
	actionMatchesAPossibleWorld bool,
) {
	worldsCache[hash] = responseWorlds
	countsCache[hash] = counts
	rewardsCache[hash] = rewards
	actionMatchesCache[hash] = actionMatchesAPossibleWorld
}

func (c *Cache) Exists(hash string) bool {
	_, ok := worldsCache[hash]
	return ok
}

func (c *Cache) GetCounts(hash string) utils.GroupedIncrementer {
	return countsCache[hash]
}

func (c *Cache) GetRewards(hash string) utils.GroupedIncrementer {
	return rewardsCache[hash]
}

func (c *Cache) GetActionMatches(hash string) bool {
	return actionMatchesCache[hash]
}

func InitUniverse(lines []string) Universe {
	universe := Universe{}
	universe.Size = len(lines)
	universe.Worlds = make([]BWorld, universe.Size)
	universe.WorldLength = len(lines[0])
	pw := NewPossibleWorlds(universe.Size)
	for index, name := range lines {
		universe.Worlds[index] = BWorld{Index: uint64(index), Name: name}
		pw.Add(uint64(index))
	}
	universe.PossibleWorlds = pw
	universe.Actions = buildAvailableActions(universe)
	return universe
}

func (u *Universe) GetWorldByIndex(index uint64) BWorld {
	return u.Worlds[index]
}

type AvailableActions struct {
	worlds       []World
	actions      []Action
	worldActions []Action
	worldSize    int
}

// TODO: Just move this to universe creation then the performance shouldn't matter
// because it'll get called once instead of every new game with NewPerspective
func buildAvailableActions(universe Universe) AvailableActions {
	aa := AvailableActions{
		actions:      []Action{},
		worldActions: make([]Action, universe.Size),
		worldSize:    universe.WorldLength,
	}

	keys := make(map[int32]bool)
	for index, world := range universe.Worlds {
		worldAction := Action(world.Name)
		aa.worldActions[index] = worldAction
		for _, runeValue := range world.Name {
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
