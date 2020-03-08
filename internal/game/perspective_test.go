package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestPerspectivePlausibleActions(t *testing.T) {
	remainingWorlds := Worlds{World("abc"), World("aab")}
	perspective := &Perspective{}
	perspective.AvailableActions = BuildAvailableActions(remainingWorlds)
	perspective.ActionFilter = TestActionFilter

	actual := perspective.PlausibleActions(remainingWorlds)
	test.AssertDeepEqual(actual, []Action{"hi"}, t)
}

func TestExploreReactions(t *testing.T) {
	worlds := Worlds{World("abc"), World("aab")}
	aa := BuildAvailableActions(worlds)
	actual := ExploreReactions(worlds, aa.AllActions(), TestAgreement)
	test.AssertDeepEqual(actual.Get("a", "a"), []World{"abc", "aab"}, t)
	test.AssertDeepEqual(actual.Get("c", "c"), []World{"abc"}, t)
	test.AssertDeepEqual(actual.Get("c", "!"), []World{"aab"}, t)
}
