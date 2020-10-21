package game

import (
	"encoding/binary"

	"github.com/coocood/freecache"
	"github.com/segmentio/fasthash/fnv1a"
)

var cacheSize = 1000 * 1024 * 1024
var cacheExpireSeconds = 0
var theFreeCache = freecache.NewCache(cacheSize)

var thePredictionsFreeCache = freecache.NewCache(cacheSize)

var PredictionCache = predictionCache{}

var TheCache = Cache{
	CacheMap: cacheMap{lookup: make(map[uint64]uint64)},
	PWCache:  pwCache{},
}

type predictionCache struct {
}

func (pc *predictionCache) Set(group uint64, value Predictions) {
	bKey := Uint64ToBytes(group)
	err := thePredictionsFreeCache.Set(bKey, value.ToBytes(), cacheExpireSeconds)
	if err != nil {
		panic(err)
	}
}

func (pc *predictionCache) Get(group uint64) Predictions {
	bKey := Uint64ToBytes(group)
	thePredictionsFreeCache.Touch(bKey, cacheExpireSeconds)
	data, err := thePredictionsFreeCache.Get(bKey)
	if err != nil {
		panic(err)
	}
	return NewPredictionsFromBytes(data)
}

func (pc *predictionCache) Exists(group uint64) bool {
	bKey := Uint64ToBytes(group)
	err := thePredictionsFreeCache.Touch(bKey, cacheExpireSeconds)
	return err == nil
}

type pwCache struct {
}

func (pwc *pwCache) Set(group uint64, value PossibleWorlds) {
	bKey := Uint64ToBytes(group)
	err := theFreeCache.Set(bKey, value.ToBytes(), cacheExpireSeconds)
	if err != nil {
		panic(err)
	}
}

func (pwc pwCache) Get(group uint64) PossibleWorlds {
	bKey := Uint64ToBytes(group)
	theFreeCache.Touch(bKey, cacheExpireSeconds)
	data, err := theFreeCache.Get(bKey)
	if err != nil {
		panic(err)
	}
	return PossibleWorldsFromBytes(data)
}

func (pwc pwCache) Exists(group uint64) bool {
	bKey := Uint64ToBytes(group)
	err := thePredictionsFreeCache.Touch(bKey, cacheExpireSeconds)
	return err == nil
}

func Uint64ToBytes(ui uint64) []byte {
	var data [8]byte
	binary.LittleEndian.PutUint64(data[:], ui)
	return data[:]
}

func BytesToUint64(data []byte) (ui uint64) {
	return binary.LittleEndian.Uint64(data[:])
}

type cacheMap struct {
	size   uint32
	lookup map[uint64]uint64
}

func (m *cacheMap) Set(group, responseHash, value uint64) {
	hash := fnv1a.AddUint64(group, responseHash)
	bKey := Uint64ToBytes(hash)
	err := theFreeCache.Set(bKey, Uint64ToBytes(value), cacheExpireSeconds)
	if err != nil {
		panic(err)
	}
}

func (m cacheMap) Get(group, responseHash uint64) uint64 {
	hash := fnv1a.AddUint64(group, responseHash)
	bKey := Uint64ToBytes(hash)
	theFreeCache.Touch(bKey, cacheExpireSeconds)
	data, err := theFreeCache.Get(bKey)
	if err != nil {
		panic(err)
	}
	return BytesToUint64(data)
}

func (m cacheMap) Exists(group, responseHash uint64) bool {
	hash := fnv1a.AddUint64(group, responseHash)
	bKey := Uint64ToBytes(hash)
	err := theFreeCache.Touch(bKey, cacheExpireSeconds)
	return err == nil
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
