package game

import "github.com/segmentio/fasthash/fnv1a"

var PredictionCache = predictionCache{lookup: make(map[uint64]uint32)}

var TheCache = Cache{
	CacheMap: cacheMap{lookup: make(map[uint64]uint64, 1000000)},
	PWCache:  pwCache{lookup: make(map[uint64]uint32)},
}

type predictionCache struct {
	values []Predictions
	size   uint32
	lookup map[uint64]uint32
}

func (pc *predictionCache) Set(group uint64, value Predictions) {
	index := pc.initializeGroup(group)
	pc.values[index] = value
}

func (pc *predictionCache) Get(group uint64) Predictions {
	index, _ := pc.lookup[group]
	return pc.values[index]
}

func (pc *predictionCache) Exists(group uint64) bool {
	_, ok := pc.lookup[group]
	return ok
}

func (pc *predictionCache) initializeGroup(group uint64) uint32 {
	pc.values = append(pc.values, Predictions{})
	pc.size++
	pc.lookup[group] = pc.size - 1
	return pc.size - 1
}

type pwCache struct {
	values []PossibleWorlds
	size   uint32
	lookup map[uint64]uint32
}

func (pwc *pwCache) Set(group uint64, value PossibleWorlds) {
	index := pwc.initializeGroup(group)
	pwc.values[index] = value
}

func (pwc pwCache) Get(group uint64) PossibleWorlds {
	index, _ := pwc.lookup[group]
	return pwc.values[index]
}

func (pwc pwCache) Exists(group uint64) bool {
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
	size   uint32
	lookup map[uint64]uint64
}

func (m *cacheMap) Set(group, responseHash, value uint64) {
	m.size++
	hash := fnv1a.AddUint64(group, responseHash)
	m.lookup[hash] = value
}

func (m cacheMap) Get(group, responseHash uint64) uint64 {
	hash := fnv1a.AddUint64(group, responseHash)
	return m.lookup[hash]
}

func (m cacheMap) Exists(group, responseHash uint64) bool {
	hash := fnv1a.AddUint64(group, responseHash)
	_, ok := m.lookup[hash]
	return ok
}

func (m cacheMap) Size() uint32 {
	return m.size
}

type Cache struct {
	CacheMap cacheMap
	PWCache  pwCache
}

func (c Cache) GetNextWorlds(gamePlayHash uint64, actionHash uint64, responseHash uint64) PossibleWorlds {
	gameActionHash := BuildGameActionHash(gamePlayHash, actionHash)
	pwHash := c.CacheMap.Get(gameActionHash, responseHash)
	return c.PWCache.Get(pwHash)
}

func BuildGameActionHash(gamePlayHash uint64, actionHash uint64) uint64 {
	return fnv1a.AddUint64(gamePlayHash, actionHash)
}

func (c *Cache) Write(
	gamePlayHash uint64,
	actionHash uint64,
	responseWorlds map[uint64]PossibleWorlds,
) {
	hash := BuildGameActionHash(gamePlayHash, actionHash)
	if exists := c.CacheMap.Exists(hash, hash); !exists {
		for responseHash, possibleWorlds := range responseWorlds {
			pwHash := possibleWorlds.SetHash()
			if exists := c.PWCache.Exists(pwHash); !exists {
				c.PWCache.Set(pwHash, possibleWorlds)
			}
			c.CacheMap.Set(hash, responseHash, pwHash)
		}
		// Hack to cache the fact that we've cached responses for this gameAction hash before
		c.CacheMap.Set(hash, hash, 0)
	}
}

func (c Cache) Exists(group, responseHash uint64) bool {
	return c.CacheMap.Exists(group, responseHash)
}

func (c Cache) Size() int {
	return int(c.CacheMap.Size())
}
