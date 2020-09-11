package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestWorldsDeprecatedSortedKey(t *testing.T) {
	worlds := BuildPossibleWorldsDeprecated([]World{
		World("ab"),
		World("co"),
		World("bi"),
	})
	expected := "abbico"
	actual := worlds.Key()
	test.AssertDeepEqual(actual, expected, t)
}

func TestWorldsSortedKey(t *testing.T) {
	worlds := BuildPossibleWorlds([]string{
		"ab",
		"co",
		"bi",
	})
	expected := "abbico"
	actual := worlds.Key()
	test.AssertDeepEqual(actual, expected, t)
}
