package utils

import (
	"fmt"
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestGroupedIncrementerGet(t *testing.T) {
	gi := NewGroupedIncrementer()
	test.AssertAlmostEqual(gi.Get("hi"), 0.0, 4, t)
	test.AssertAlmostEqual(gi.Get("bye"), 0.0, 4, t)

	gi.Increment("hi", 1.1)
	test.AssertAlmostEqual(gi.Get("hi"), 1.1, 4, t)
}

func TestGroupedIncrementerTotal(t *testing.T) {
	gi := NewGroupedIncrementer()
	test.AssertAlmostEqual(gi.Total(), 0.0, 4, t)

	gi.Increment("hi", 1.1)
	gi.Increment("hi", 1.2)
	gi.Increment("bye", 1.4)

	test.AssertAlmostEqual(gi.Total(), 3.7, 4, t)
}

func BenchmarkTotalSmall(b *testing.B) {
	gi := NewGroupedIncrementer()
	for i := 0; i < 10; i++ {
		gi.Increment(fmt.Sprintf("%d", i), 2.4)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gi.Total()
	}
}

func BenchmarkTotalLarge(b *testing.B) {
	gi := NewGroupedIncrementer()
	for i := 0; i < 1000000; i++ {
		gi.Increment(fmt.Sprintf("%d", i), 2.4)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gi.Total()
	}
}

func TestGroupedIncrementerAll(t *testing.T) {
	gi := NewGroupedIncrementer()
	test.AssertDeepEqual(gi.All(), []float64{}, t)

	gi.Increment("hi", 1.1)
	gi.Increment("hi", 1.2)
	gi.Increment("bye", 1.4)
	test.AssertDeepEqual(gi.All(), []float64{2.3, 1.4}, t)
}

func BenchmarkAllSmall(b *testing.B) {
	gi := NewGroupedIncrementer()
	numElements := 10
	for i := 0; i < numElements; i++ {
		gi.Increment(fmt.Sprintf("%d", i), 2.4)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gi.All()
	}
}

func BenchmarkAllLarge(b *testing.B) {
	gi := NewGroupedIncrementer()
	numElements := 1000000
	for i := 0; i < numElements; i++ {
		gi.Increment(fmt.Sprintf("%d", i), 2.4)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gi.All()
	}
}
