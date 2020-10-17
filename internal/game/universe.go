package game

import (
	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/segmentio/fasthash/fnv1a"
)

type Response string
type Act func(gp GamePlay, w World, action Action) (reaction Reaction)
type GetActions func(gp GamePlay, w PossibleWorlds, aa AvailableActions) []Action

func NewWorld(index uint64, worldString string) World {
	w := World{}
	w.Index = index
	w.Name = worldString
	return w
}

type Action struct {
	Bytes  []byte
	Hash   uint64
	Length int
}

func (a Action) String() string {
	return string(a.Bytes)
}

func (a Action) Equals(other Action) bool {
	return a.Hash == other.Hash
}

func NewAction(actionString string) Action {
	a := Action{}
	a.Bytes = []byte(actionString)
	a.Hash = fnv1a.HashBytes64(a.Bytes)
	a.Length = len(actionString)
	return a
}

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
	Reactions   []Reaction
	TurnCount   int
	TotalReward float64
	TotalPain   float64
	Hash        uint64
}

func InitGamePlay() GamePlay {
	gp := GamePlay{}
	gp.Hash = fnv1a.Init64
	return gp
}

func (gp *GamePlay) AddTurn(reaction Reaction) {
	gp.Reactions = append(gp.Reactions, reaction)
	gp.TurnCount++
	gp.TotalReward += reaction.Reward
	gp.TotalPain += 1.0 - reaction.Reward
	gp.Hash = fnv1a.AddBytes64(gp.Hash, reaction.Action.Bytes)
}

func (gp *GamePlay) PreviouslyActed(action Action) bool {
	for _, reaction := range gp.Reactions {
		if reaction.Action.Equals(action) {
			return true
		}
	}
	return false
}

type PossibleWorlds struct {
	Bits     roaring64.Bitmap
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
	pw.Bits = *roaring64.BitmapOf(length64 + 1)
	pw.count = 0
	pw.Capacity = length
	return pw
}

func BuildActionWorldHash(action Action, worldsHash uint64) uint64 {
	worldsHash = fnv1a.AddBytes64(worldsHash, action.Bytes)
	return worldsHash

}

type World struct {
	Index uint64
	Name  string
}

type Universe struct {
	Worlds         []World
	WorldLength    int
	Actions        AvailableActions
	PossibleWorlds PossibleWorlds
	Size           int
}

var BUniverse = Universe{}

type Cache struct{}

var TheCache = Cache{}

var thePWCache = pwCache{lookup: make(map[uint64]uint32)}
var theCacheMap = cacheMap{lookup: make(map[uint64]uint32)}

type pwCache struct {
	values []PossibleWorlds
	size   uint32
	lookup map[uint64]uint32
}

func (pwc *pwCache) Set(group uint64, value PossibleWorlds) {
	index := pwc.initializeGroup(group)
	pwc.values[index] = value
}

func (pwc *pwCache) Get(group uint64) PossibleWorlds {
	index, _ := pwc.lookup[group]
	return pwc.values[index]
}

func (pwc *pwCache) Exists(group uint64) bool {
	_, ok := pwc.lookup[group]
	return ok
}

func (pwc *pwCache) initializeGroup(group uint64) uint32 {
	pwc.values = append(pwc.values, PossibleWorlds{})
	pwc.size++
	pwc.lookup[group] = pwc.size - 1
	return pwc.size - 1
}

type cacheMap struct {
	values []ActionCacheModel
	size   uint32
	lookup map[uint64]uint32
}

func (m *cacheMap) Set(group uint64, value ActionCacheModel) {
	index := m.initializeGroup(group)
	m.values[index] = value
}

func (m *cacheMap) Get(group uint64) ActionCacheModel {
	index, _ := m.lookup[group]
	return m.values[index]
}

func (m *cacheMap) Exists(group uint64) bool {
	_, ok := m.lookup[group]
	return ok
}

func (m *cacheMap) Size() uint32 {
	return m.size
}

func (m *cacheMap) initializeGroup(group uint64) uint32 {
	m.values = append(m.values, ActionCacheModel{})
	m.size++
	m.lookup[group] = m.size - 1
	return m.size - 1
}

type ActionCacheModel struct {
	Worlds map[uint64]uint64
}

func (c *Cache) GetNextWorlds(hash uint64, response string) PossibleWorlds {
	responseHash := fnv1a.HashString64(response)
	pwHash := theCacheMap.Get(hash).Worlds[responseHash]
	return thePWCache.Get(pwHash)
}

func (c *Cache) Write(
	hash uint64,
	responseWorlds map[uint64]PossibleWorlds,
) {
	if exists := theCacheMap.Exists(hash); !exists {
		responseCache := make(map[uint64]uint64)
		for responseHash, possibleWorlds := range responseWorlds {
			pwHash := possibleWorlds.SetHash()
			if exists := thePWCache.Exists(pwHash); !exists {
				thePWCache.Set(pwHash, possibleWorlds)
			}
			responseCache[responseHash] = pwHash
		}
		theCacheMap.Set(hash, ActionCacheModel{
			Worlds: responseCache,
		})
	}
}

func (c *Cache) Exists(hash uint64) bool {
	return theCacheMap.Exists(hash)
}

func (c *Cache) Size() int {
	return int(theCacheMap.Size())
}

func InitUniverse(lines []string) Universe {
	universe := Universe{}
	universe.Size = len(lines)
	universe.Worlds = make([]World, universe.Size)
	universe.WorldLength = len(lines[0])
	pw := NewPossibleWorlds(universe.Size)
	for index, name := range lines {
		universe.Worlds[index] = World{Index: uint64(index), Name: name}
		pw.Add(uint64(index))
	}
	pw.SetHash()
	universe.PossibleWorlds = pw
	universe.Actions = buildAvailableActions(universe)
	return universe
}

func (u *Universe) GetWorldByIndex(index uint64) World {
	return u.Worlds[index]
}

type AvailableActions struct {
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
		worldAction := NewAction(world.Name)
		aa.worldActions[index] = worldAction
		for _, runeValue := range world.Name {
			if _, ok := keys[runeValue]; !ok {
				keys[runeValue] = true
				aa.actions = append(aa.actions, NewAction(string(runeValue)))
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
