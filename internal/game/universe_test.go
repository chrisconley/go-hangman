package game

import (
	"fmt"
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestPossibleWorldsGetHash(t *testing.T) {
	pw1 := NewPossibleWorlds(100)
	pw1.Add(1)
	pw1.Add(99)

	pw2 := NewPossibleWorlds(100)
	pw2.Add(2)
	pw2.Add(98)
	test.AssertDeepNotEqual(pw1.GetHash(), pw2.GetHash(), t)
}

func TestHashingConflicts(t *testing.T) {
	pw1 := NewPossibleWorlds(25011)
	pw1.Add(uint64(617))
	pw1.Add(uint64(618))

	pw2 := NewPossibleWorlds(25011)
	pw2.Add(uint64(809))
	pw2.Add(uint64(810))

	action := Action("t")

	test.AssertDeepNotEqual(
		BuildActionWorldHash(action, pw1.GetHash()),
		BuildActionWorldHash(action, pw2.GetHash()),
		t,
	)
}

func TestAvailableActions(t *testing.T) {
	InitUniverse([]string{"aa", "bb"})
	actual := BUniverse.Actions.AllActions()
	expected := []Action{"a", "b", "aa", "bb"}
	test.AssertDeepEqual(actual, expected, t)
}

func TestAvailableActionsWorldActionsLength(t *testing.T) {
	InitUniverse([]string{"aa", "bb"})
	actual := BUniverse.Actions.WorldActionsLength()
	test.AssertDeepEqual(actual, 2, t)
}

func TestAvailableActionsActionsLength(t *testing.T) {
	InitUniverse([]string{"aa", "bb"})
	actual := BUniverse.Actions.ActionsLength()
	test.AssertDeepEqual(actual, 2, t)
}

func TestAvailableActionsActions(t *testing.T) {
	InitUniverse([]string{"aa", "bb"})
	actual := BUniverse.Actions.Actions()
	test.AssertDeepEqual(actual, []Action{"a", "b"}, t)
}

func TestAvailableActionsWorldActions(t *testing.T) {
	InitUniverse([]string{"aa", "bb"})
	actual := BUniverse.Actions.WorldActions()
	test.AssertDeepEqual(actual, []Action{"aa", "bb"}, t)
}

func BenchmarkAvailableActionsActionsLengthSmall(b *testing.B) {
	worlds := []string{}
	numElements := 10
	for i := 0; i < numElements; i++ {
		worlds = append(worlds, fmt.Sprintf("hi%d", i))
	}
	InitUniverse(worlds)
	aa := BUniverse.Actions
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
	InitUniverse(worlds)
	aa := BUniverse.Actions
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
