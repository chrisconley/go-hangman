package game

import (
	"encoding/binary"
	"time"

	"github.com/allegro/bigcache"
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
var thePredictionsFreeCache, _ = bigcache.NewBigCache(bigCacheConfig)

var PredictionCache = predictionCache{}

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

func Uint64ToBytes(ui uint64) []byte {
	var data [8]byte
	binary.LittleEndian.PutUint64(data[:], ui)
	return data[:]
}

func BytesToUint64(data []byte) (ui uint64) {
	return binary.LittleEndian.Uint64(data[:])
}
