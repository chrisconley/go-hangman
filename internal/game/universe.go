package game

import (
	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/segmentio/fasthash/fnv1a"
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

type GamePlay struct {
	Turns     []Reaction
	TurnCount int
}

type PossibleWorlds struct {
	Bits     *roaring64.Bitmap
	count    int
	Capacity int
	hash     uint64
}

func (pw *PossibleWorlds) GetHash() uint64 {
	return pw.hash
}

func (pw *PossibleWorlds) SetHash() uint64 {
	bytes, _ := pw.Bits.ToBytes()
	pw.hash = fnv1a.Init64
	pw.hash = fnv1a.AddBytes64(pw.hash, bytes)
	return pw.hash
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

func BuildActionWorldHash(action Action, worldsHash uint64) uint64 {
	worldsHash = fnv1a.AddString64(worldsHash, string(action))
	return worldsHash

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
var possibleWorldsCache = make(map[uint64]PossibleWorlds)
var theCache = make(map[uint64]ActionCacheModel)

type ActionCacheModel struct {
	Worlds        map[string]uint64
	Counts        utils.GroupedIncrementer
	Rewards       utils.GroupedIncrementer
	ActionMatches bool
}

func (c *Cache) GetNextWorlds(hash uint64, response string) PossibleWorlds {
	pwHash := theCache[hash].Worlds[response]
	return possibleWorldsCache[pwHash]
}

func (c *Cache) Write(
	hash uint64,
	responseWorlds map[string]PossibleWorlds,
	counts utils.GroupedIncrementer,
	rewards utils.GroupedIncrementer,
	actionMatchesAPossibleWorld bool,
) {
	responseCache := make(map[string]uint64)
	for response, possibleWorlds := range responseWorlds {
		pwHash := possibleWorlds.SetHash()
		possibleWorldsCache[pwHash] = possibleWorlds
		responseCache[response] = pwHash
	}
	theCache[hash] = ActionCacheModel{
		Worlds:        responseCache,
		Counts:        counts,
		Rewards:       rewards,
		ActionMatches: actionMatchesAPossibleWorld,
	}
}

func (c *Cache) Exists(hash uint64) bool {
	_, ok := theCache[hash]
	return ok
}

func (c *Cache) GetCounts(hash uint64) utils.GroupedIncrementer {
	return theCache[hash].Counts
}

func (c *Cache) GetRewards(hash uint64) utils.GroupedIncrementer {
	return theCache[hash].Rewards
}

func (c *Cache) GetActionMatches(hash uint64) bool {
	return theCache[hash].ActionMatches
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
	pw.SetHash()
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
