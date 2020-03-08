package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestExploreReactions(t *testing.T) {
	worlds := Worlds{World("abc"), World("aab")}
	aa := BuildAvailableActions(worlds)
	actual := ExploreReactions(worlds, aa, TestAgreement)
	test.AssertDeepEqual(actual.Get("a", "a"), []World{"abc", "aab"}, t)
	test.AssertDeepEqual(actual.Get("c", "c"), []World{"abc"}, t)
	test.AssertDeepEqual(actual.Get("c", "!"), []World{"aab"}, t)
}
