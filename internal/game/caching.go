package game

import "github.com/segmentio/fasthash/fnv1a"

var PredictionCache = make(map[uint64]Predictions)

var TheCache = Cache{}

var theCacheMap = cacheMap{lookup: make(map[GameActionHash]map[uint64]uint64)}

var thePWCache = pwCache{lookup: make(map[uint64]uint32)}

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

type Hash interface{}
type hash uint64

type GameActionHash interface {
	Hash
}
type gameActionHash hash

type cacheMap struct {
	size   uint32
	lookup map[GameActionHash]map[uint64]uint64
}

func (m *cacheMap) Set(group GameActionHash, value map[uint64]uint64) {
	m.size++
	m.lookup[group] = value
}

func (m *cacheMap) Get(group GameActionHash) map[uint64]uint64 {
	responses := m.lookup[group]
	return responses
}

func (m *cacheMap) Exists(group GameActionHash) bool {
	_, ok := m.lookup[group]
	return ok
}

func (m *cacheMap) Size() uint32 {
	return m.size
}

type Cache struct{}

func (c *Cache) GetNextWorlds(gamePlayHash uint64, actionHash uint64, responseHash uint64) PossibleWorlds {
	gameActionHash := BuildGameActionHash(gamePlayHash, actionHash)
	pwHash := theCacheMap.Get(gameActionHash)[responseHash]
	return thePWCache.Get(pwHash)
}

func BuildGameActionHash(gamePlayHash uint64, actionHash uint64) GameActionHash {
	return fnv1a.AddUint64(gamePlayHash, actionHash)
}

func (c *Cache) Write(
	gamePlayHash uint64,
	actionHash uint64,
	responseWorlds map[uint64]PossibleWorlds,
) {
	hash := BuildGameActionHash(gamePlayHash, actionHash)
	if exists := theCacheMap.Exists(hash); !exists {
		responseCache := make(map[uint64]uint64)
		for responseHash, possibleWorlds := range responseWorlds {
			pwHash := possibleWorlds.SetHash()
			if exists := thePWCache.Exists(pwHash); !exists {
				thePWCache.Set(pwHash, possibleWorlds)
			}
			responseCache[responseHash] = pwHash
		}
		theCacheMap.Set(hash, responseCache)
	}
}

func (c *Cache) Exists(hash uint64) bool {
	return theCacheMap.Exists(hash)
}

func (c *Cache) Size() int {
	return int(theCacheMap.Size())
}
