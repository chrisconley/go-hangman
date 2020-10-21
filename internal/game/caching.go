package game

import (
	"github.com/coocood/freecache"
	"github.com/segmentio/fasthash/fnv1a"
)

var PredictionCache = predictionCache{}

var TheCache = Cache{
	CacheMap: cacheMap{lookup: make(map[uint64]uint64)},
	PWCache:  pwCache{},
}

type predictionCache struct {
}

func (pc *predictionCache) Set(group uint64, value Predictions) {
	err := theFreeCache.SetInt(int64(group), value.ToBytes(), 0)
	if err != nil {
		panic(err)
	}
}

func (pc *predictionCache) Get(group uint64) Predictions {
	data, err := theFreeCache.GetInt(int64(group))
	if err != nil {
		panic(err)
	}
	return NewPredictionsFromBytes(data)
}

func (pc *predictionCache) Exists(group uint64) bool {
	_, err := theFreeCache.GetInt(int64(group))
	return err == nil
}

var cacheSize = 1000 * 1024 * 1024
var theFreeCache = freecache.NewCache(cacheSize)

type pwCache struct {
}

func (pwc *pwCache) Set(group uint64, value PossibleWorlds) {
	err := theFreeCache.SetInt(int64(group), value.ToBytes(), 0)
	if err != nil {
		panic(err)
	}
}

func (pwc pwCache) Get(group uint64) PossibleWorlds {
	data, err := theFreeCache.GetInt(int64(group))
	if err != nil {
		panic(err)
	}
	return PossibleWorldsFromBytes(data)
}

func (pwc pwCache) Exists(group uint64) bool {
	_, err := theFreeCache.GetInt(int64(group))
	return err == nil
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
