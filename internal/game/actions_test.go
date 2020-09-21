package game

import (
	"fmt"
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestBuildAvailableActions(t *testing.T) {
	worlds := BuildPossibleWorlds([]string{"aa", "bb"})
	aa := BuildAvailableActions(worlds)
	actual := aa.AllActions()
	expected := []Action{"a", "b", "aa", "bb"}
	test.AssertDeepEqual(actual, expected, t)
}

func TestAvailableActionsWorldActionsLength(t *testing.T) {
	worlds := BuildPossibleWorlds([]string{"aa", "bb"})
	aa := BuildAvailableActions(worlds)
	actual := aa.WorldActionsLength()
	test.AssertDeepEqual(actual, 2, t)
}

func TestAvailableActionsActionsLength(t *testing.T) {
	worlds := BuildPossibleWorlds([]string{"aa", "bb"})
	aa := BuildAvailableActions(worlds)
	actual := aa.ActionsLength()
	test.AssertDeepEqual(actual, 2, t)
}

func BenchmarkAvailableActionsActionsLengthSmall(b *testing.B) {
	worlds := []string{}
	numElements := 10
	for i := 0; i < numElements; i++ {
		worlds = append(worlds, fmt.Sprintf("hi%d", i))
	}
	aa := BuildAvailableActions(BuildPossibleWorlds(worlds))
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
	aa := BuildAvailableActions(BuildPossibleWorlds(worlds))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aa.ActionsLength()
	}
}
