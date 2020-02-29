package utils

type GroupedIncrementer struct {
	values []float64
	lookup map[string]int
	total  float64
}

func NewGroupedIncrementer() GroupedIncrementer {
	gi := GroupedIncrementer{}
	gi.values = []float64{}
	gi.lookup = make(map[string]int)
	gi.total = 0.0
	return gi
}

func (gi *GroupedIncrementer) Increment(group string, value float64) {
	gi.findOrInitializeGroup(group)
	index := gi.lookup[group]
	gi.values[index] += value
	gi.total += value
}

func (gi *GroupedIncrementer) Get(group string) float64 {
	gi.findOrInitializeGroup(group)
	index := gi.lookup[group]
	return gi.values[index]
}

func (gi *GroupedIncrementer) Total() float64 {
	return gi.total
}

func (gi *GroupedIncrementer) All() []float64 {
	return gi.values
}

func (gi *GroupedIncrementer) findOrInitializeGroup(group string) {
	_, ok := gi.lookup[group]
	if !ok {
		gi.values = append(gi.values, 0.0)
		length := len(gi.values)
		gi.lookup[group] = length - 1
	}
}
