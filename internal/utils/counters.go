package utils

import "github.com/segmentio/fasthash/fnv1a"

type GroupedIncrementer struct {
	values []float64
	size   uint16
	lookup map[uint64]uint16
	total  float64
}

func NewGroupedIncrementer() GroupedIncrementer {
	gi := GroupedIncrementer{}
	gi.values = []float64{}
	gi.lookup = make(map[uint64]uint16)
	gi.total = 0.0
	gi.size = 0
	return gi
}

func (gi *GroupedIncrementer) Increment(group string, value float64) {
	index := gi.findOrInitializeGroup(group)
	gi.values[index] += value
	gi.total += value
}

func (gi *GroupedIncrementer) Get(group string) float64 {
	index := gi.findOrInitializeGroup(group)
	return gi.values[index]
}

func (gi *GroupedIncrementer) Total() float64 {
	return gi.total
}

func (gi *GroupedIncrementer) All() []float64 {
	return gi.values
}

func (gi *GroupedIncrementer) GetHash(groups []string) uint64 {
	hash := fnv1a.Init64
	for _, group := range groups {
		index := gi.findOrInitializeGroup(group)
		count := gi.values[index]
		hash = fnv1a.AddString64(hash, group)
		hash = fnv1a.AddUint64(hash, uint64(count))
	}
	return hash
}

func (gi *GroupedIncrementer) findOrInitializeGroup(group string) uint16 {
	hash := fnv1a.Init64
	hash = fnv1a.AddString64(hash, group)
	index, ok := gi.lookup[hash]
	if !ok {
		gi.values = append(gi.values, 0.0)
		gi.size++
		gi.lookup[hash] = gi.size - 1
		return gi.size - 1
	}
	return index

}
