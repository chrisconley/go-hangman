package game

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/coocood/freecache"
)

// The default cache size is  100mb
var cacheSize = 100 * 1024 * 1024
var cacheExpireSeconds = 5 * 60
var ThePredictionsFreeCache = freecache.NewCache(cacheSize)

var PredictionCache = predictionCache{}

type predictionCache struct {
	size uint64
}

func (pc *predictionCache) Set(group uint64, value Predictions) {
	bKey := Uint64ToBytes(group)
	err := ThePredictionsFreeCache.Set(bKey, value.ToBytes(), cacheExpireSeconds)
	if err != nil {
		fmt.Println(value.PredictionsMap.Size(), len(value.ToBytes()))
		panic(err)
	}
	pc.size++
}

func (pc *predictionCache) Get(group uint64) (predictions Predictions, err error) {
	bKey := Uint64ToBytes(group)
	data, err := ThePredictionsFreeCache.Get(bKey)
	if err != nil {
		predictions = NewPredictions()
	} else {
		predictions = NewPredictionsFromBytes(data)
	}
	return predictions, err
}

type freecacheStats struct {
	EntryCount    int64
	EvacuateCount int64
	ExpiredCount  int64
	HitCount      int64
	MissCount     int64
	TouchedCount  int64
}

func FreeCacheStats(cache *freecache.Cache) freecacheStats {
	return freecacheStats{
		EntryCount:    cache.EntryCount(),
		EvacuateCount: cache.EvacuateCount(),
		ExpiredCount:  cache.ExpiredCount(),
		HitCount:      cache.HitCount(),
		MissCount:     cache.MissCount(),
		TouchedCount:  cache.TouchedCount(),
	}
}

func (pc *predictionCache) Size() uint64 {
	return pc.size
}

var ObserveCache = observeCache{
	Cache: freecache.NewCache(cacheSize),
}

type observeCache struct {
	Cache *freecache.Cache
	size  uint64
}

func (oc *observeCache) Set(hash uint64, possibleWorlds PossibleWorlds) {
	err := oc.Cache.SetInt(int64(hash), possibleWorlds.ToBytes(), 5*60)
	if err != nil {
		panic(err)
	}
	oc.size++
}

func (oc *observeCache) Get(hash uint64) (possibleWorlds PossibleWorlds, err error) {
	data, err := oc.Cache.GetInt(int64(hash))
	if err != nil {
		possibleWorlds = NewPossibleWorlds(BUniverse.Size)
	} else {
		possibleWorlds = PossibleWorldsFromBytes(data)
	}
	return possibleWorlds, err
}

func Uint64ToBytes(ui uint64) []byte {
	var data [8]byte
	binary.LittleEndian.PutUint64(data[:], ui)
	return data[:]
}

func BytesToUint64(data []byte) (ui uint64) {
	return binary.LittleEndian.Uint64(data[:])
}

func Float64ToBytes(f float64) []byte {
	return Uint64ToBytes(math.Float64bits(f))
}

func BytesToFloat64(data []byte) float64 {
	return math.Float64frombits(BytesToUint64(data[:]))
}
