package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestBuildAvailableActions(t *testing.T) {
	aa := BuildAvailableActions([]World{"aa", "bb"})
	actual := aa.All()
	expected := []Action{"aa", "a", "bb", "b"}
	test.AssertDeepEqual(actual, expected, t)
}

func TestAvailableActionsUpdate(t *testing.T) {
	aa := BuildAvailableActions([]World{"aa", "bb"})
	aa.Update(Action("a"))
	actual := aa.All()
	expected := []Action{"aa", "bb", "b"}
	test.AssertDeepEqual(actual, expected, t)
}
