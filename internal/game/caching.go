package game

import "github.com/segmentio/fasthash/fnv1a"

var PredictionCache = make(map[uint64]Predictions)

var TheCache = Cache{}

var theCacheMap = cacheMap{lookup: make(map[uint64]uint32)}

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

type ActionCacheModel struct {
	Worlds map[uint64]uint64
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

type Cache struct{}

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
