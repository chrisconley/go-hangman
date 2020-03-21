package game

import (
	"fmt"
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestBuildAvailableActions(t *testing.T) {
	worlds := BuildPossibleWorlds([]World{"aa", "bb"})
	aa := BuildAvailableActions(worlds)
	actual := aa.AllActions()
	expected := []Action{"a", "b", "aa", "bb"}
	test.AssertDeepEqual(actual, expected, t)
}

func TestAvailableActionsUpdate(t *testing.T) {
	worlds := BuildPossibleWorlds([]World{"aa", "bb"})
	aa := BuildAvailableActions(worlds)
	aa.Update(Action("a"))
	actual := aa.AllActions()
	expected := []Action{"b", "aa", "bb"}
	test.AssertDeepEqual(actual, expected, t)
}

func TestAvailableActionsWorldActionsLength(t *testing.T) {
	worlds := BuildPossibleWorlds([]World{"aa", "bb"})
	aa := BuildAvailableActions(worlds)
	actual := aa.WorldActionsLength()
	test.AssertDeepEqual(actual, 2, t)

	aa.Update(Action("a"))
	actual = aa.WorldActionsLength()
	test.AssertDeepEqual(actual, 2, t)

	aa.Update(Action("aa"))
	actual = aa.WorldActionsLength()
	test.AssertDeepEqual(actual, 1, t)
}

func TestAvailableActionsActionsLength(t *testing.T) {
	worlds := BuildPossibleWorlds([]World{"aa", "bb"})
	aa := BuildAvailableActions(worlds)
	actual := aa.ActionsLength()
	test.AssertDeepEqual(actual, 2, t)

	aa.Update(Action("a"))
	actual = aa.ActionsLength()
	test.AssertDeepEqual(actual, 1, t)

	aa.Update(Action("aa"))
	actual = aa.ActionsLength()
	test.AssertDeepEqual(actual, 1, t)
}

func BenchmarkAvailableActionsActionsLengthSmall(b *testing.B) {
	worlds := []World{}
	numElements := 10
	for i := 0; i < numElements; i++ {
		worlds = append(worlds, World(fmt.Sprintf("hi%d", i)))
	}
	aa := BuildAvailableActions(BuildPossibleWorlds(worlds))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aa.ActionsLength()
	}
}

func BenchmarkAvailableActionsActionsLengthLarge(b *testing.B) {
	worlds := []World{}
	numElements := 10000
	for i := 0; i < numElements; i++ {
		worlds = append(worlds, World(fmt.Sprintf("hi%d", i)))
	}
	aa := BuildAvailableActions(BuildPossibleWorlds(worlds))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aa.ActionsLength()
	}
}
