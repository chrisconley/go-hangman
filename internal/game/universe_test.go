package game

import (
	"fmt"
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestAvailableActions(t *testing.T) {
	universe := InitUniverse([]string{"aa", "bb"})
	actual := universe.Actions.AllActions()
	expected := []Action{
		NewAction("a"),
		NewAction("b"),
		NewAction("aa"),
		NewAction("bb"),
	}
	test.AssertDeepEqual(actual, expected, t)
}

func TestAvailableActionsWorldActionsLength(t *testing.T) {
	universe := InitUniverse([]string{"aa", "bb"})
	actual := universe.Actions.WorldActionsLength()
	test.AssertDeepEqual(actual, 2, t)
}

func TestAvailableActionsActionsLength(t *testing.T) {
	universe := InitUniverse([]string{"aa", "bb"})
	actual := universe.Actions.ActionsLength()
	test.AssertDeepEqual(actual, 2, t)
}

func TestAvailableActionsActions(t *testing.T) {
	universe := InitUniverse([]string{"aa", "bb"})
	actual := universe.Actions.Actions()
	test.AssertDeepEqual(actual, []Action{NewAction("a"), NewAction("b")}, t)
}

func TestAvailableActionsWorldActions(t *testing.T) {
	universe := InitUniverse([]string{"aa", "bb"})
	actual := universe.Actions.WorldActions()
	test.AssertDeepEqual(actual, []Action{NewAction("aa"), NewAction("bb")}, t)
}

func BenchmarkAvailableActionsActionsLengthSmall(b *testing.B) {
	worlds := []string{}
	numElements := 10
	for i := 0; i < numElements; i++ {
		worlds = append(worlds, fmt.Sprintf("hi%d", i))
	}
	universe := InitUniverse(worlds)
	aa := universe.Actions
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aa.ActionsLength()
	}
}

func BenchmarkAvailableActionsActionsLengthLarge(b *testing.B) {
	worlds := []string{}
	numElements := 10000
	for i := 0; i < numElements; i++ {
		worlds = append(worlds, fmt.Sprintf("hi%d", i))
	}
	universe := InitUniverse(worlds)
	aa := universe.Actions
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aa.ActionsLength()
	}
}

func TestPossibleWorldsLen(t *testing.T) {
	pw := NewPossibleWorlds(10)
	test.AssertDeepEqual(pw.Len(), 0, t)
	pw.Add(2)
	test.AssertDeepEqual(pw.Len(), 1, t)
	pw.Add(5)
	pw.Add(6)
	test.AssertDeepEqual(pw.Len(), 3, t)

}

func TestPossibleWorldsGetHash(t *testing.T) {
	pw1 := NewPossibleWorlds(100)
	pw1.Add(1)
	pw1.Add(99)
	pw1.SetHash()

	pw2 := NewPossibleWorlds(100)
	pw2.Add(2)
	pw2.Add(98)
	pw2.SetHash()
	test.AssertDeepNotEqual(pw1.GetHash(), pw2.GetHash(), t)
}

func TestPossibleWorldsMarshalling(t *testing.T) {
	pw1 := NewPossibleWorlds(100)
	pw1.Add(1)
	pw1.Add(99)
	pw1.SetHash()

	pw1Bytes := pw1.ToBytes()
	fmt.Println("bytes length", len(pw1Bytes))

	pw2 := PossibleWorldsFromBytes(pw1Bytes)

	test.AssertDeepEqual(pw1.Bits.ToArray(), pw2.Bits.ToArray(), t)
	test.AssertDeepEqual(pw1.Len(), pw2.Len(), t)
	test.AssertDeepEqual(pw1.Capacity, pw2.Capacity, t)
	test.AssertDeepEqual(pw1.GetHash(), pw2.GetHash(), t)

}
