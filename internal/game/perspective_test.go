package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestExploreReactions(t *testing.T) {
	worlds := PossibleWorlds{PossibleWorld{"abc", 0.5}, PossibleWorld{"aab", 0.5}}
	actions := PossibleActions{"a", "b", "c"}
	actual := ExploreReactions(worlds, actions, TestAgreement)
	test.AssertDeepEqual(actual.Get("a", "a"), []World{"abc", "aab"}, t)
	test.AssertDeepEqual(actual.Get("c", "c"), []World{"abc"}, t)
	test.AssertDeepEqual(actual.Get("c", "!"), []World{"aab"}, t)
}
