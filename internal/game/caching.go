package game

import (
	"encoding/binary"
	"time"

	"github.com/allegro/bigcache"
	"github.com/segmentio/fasthash/fnv1a"
)

var cacheSize = 1000 * 1024 * 1024
var cacheExpireSeconds = 0

// var theFreeCache = freecache.NewCache(cacheSize)
// var thePredictionsFreeCache = freecache.NewCache(cacheSize)
var bigCacheConfig = bigcache.Config{
	// number of shards (must be a power of 2)
	Shards: 1024,

	// time after which entry can be evicted
	LifeWindow: 10 * time.Hour,

	// Interval between removing expired entries (clean up).
	// If set to <= 0 then no action is performed.
	// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
	CleanWindow: 0,

	// rps * lifeWindow, used only in initial memory allocation
	MaxEntriesInWindow: 100 * 10 * 60,

	// max entry size in bytes, used only in initial memory allocation
	MaxEntrySize: 10000,

	// prints information about additional memory allocation
	Verbose: false,

	// cache will not allocate more memory than this limit, value in MB
	// if value is reached then the oldest entries can be overridden for the new ones
	// 0 value means no size limit
	HardMaxCacheSize: 0,

	// callback fired when the oldest entry is removed because of its expiration time or no space left
	// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
	// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
	OnRemove: nil,

	// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
	// for the new entry, or because delete was called. A constant representing the reason will be passed through.
	// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
	// Ignored if OnRemove is specified.
	OnRemoveWithReason: nil,
}
var theFreeCache, _ = bigcache.NewBigCache(bigCacheConfig)
var thePredictionsFreeCache, _ = bigcache.NewBigCache(bigCacheConfig)

var PredictionCache = predictionCache{}

var TheCache = Cache{
	CacheMap: cacheMap{lookup: make(map[uint64]uint64)},
	PWCache:  pwCache{},
}

type predictionCache struct {
}

func (pc *predictionCache) Set(group uint64, value Predictions) {
	bKey := Uint64ToBytes(group)
	// err := thePredictionsFreeCache.Set(bKey, value.ToBytes(), cacheExpireSeconds)
	err := thePredictionsFreeCache.Set(string(bKey), value.ToBytes())
	if err != nil {
		panic(err)
	}
}

func (pc *predictionCache) Get(group uint64) Predictions {
	bKey := Uint64ToBytes(group)
	// thePredictionsFreeCache.Touch(bKey, cacheExpireSeconds)
	data, err := thePredictionsFreeCache.Get(string(bKey))
	if err != nil {
		panic(err)
	}
	return NewPredictionsFromBytes(data)
}

func (pc *predictionCache) Exists(group uint64) bool {
	bKey := Uint64ToBytes(group)
	// err := thePredictionsFreeCache.Touch(bKey, cacheExpireSeconds)
	_, err := thePredictionsFreeCache.Get(string(bKey))
	return err == nil
}

type pwCache struct {
}

func (pwc *pwCache) Set(group uint64, value PossibleWorlds) {
	bKey := Uint64ToBytes(group)
	// err := theFreeCache.Set(bKey, value.ToBytes(), cacheExpireSeconds)
	err := theFreeCache.Set(string(bKey), value.ToBytes())
	if err != nil {
		panic(err)
	}
}

func (pwc pwCache) Get(group uint64) PossibleWorlds {
	bKey := Uint64ToBytes(group)
	// theFreeCache.Touch(bKey, cacheExpireSeconds)
	data, err := theFreeCache.Get(string(bKey))
	if err != nil {
		panic(err)
	}
	return PossibleWorldsFromBytes(data)
}

func (pwc pwCache) Exists(group uint64) bool {
	bKey := Uint64ToBytes(group)
	// err := thePredictionsFreeCache.Touch(bKey, cacheExpireSeconds)
	_, err := theFreeCache.Get(string(bKey))
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
	// err := theFreeCache.Set(bKey, Uint64ToBytes(value), cacheExpireSeconds)
	err := theFreeCache.Set(string(bKey), Uint64ToBytes(value))
	if err != nil {
		panic(err)
	}
}

func (m cacheMap) Get(group, responseHash uint64) uint64 {
	hash := fnv1a.AddUint64(group, responseHash)
	bKey := Uint64ToBytes(hash)
	// theFreeCache.Touch(bKey, cacheExpireSeconds)
	// data, err := theFreeCache.Get(bKey)
	data, err := theFreeCache.Get(string(bKey))
	if err != nil {
		panic(err)
	}
	return BytesToUint64(data)
}

func (m cacheMap) Exists(group, responseHash uint64) bool {
	hash := fnv1a.AddUint64(group, responseHash)
	bKey := Uint64ToBytes(hash)
	// err := theFreeCache.Touch(bKey, cacheExpireSeconds)
	_, err := theFreeCache.Get(string(bKey))
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

func (c Cache) Exists(gamePlayHash uint64, actionHash uint64, responseHash uint64) bool {
	hash := BuildGameActionHash(gamePlayHash, actionHash)
	return c.CacheMap.Exists(hash, responseHash)
}

func (c Cache) Size() int {
	return int(c.CacheMap.Size())
}
