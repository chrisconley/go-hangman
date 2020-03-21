package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestWorldsSortedKey(t *testing.T) {
	worlds := BuildPossibleWorlds([]World{
		World("ab"),
		World("co"),
		World("bi"),
	})
	expected := "abbico"
	actual := worlds.Key()
	test.AssertDeepEqual(actual, expected, t)
}
