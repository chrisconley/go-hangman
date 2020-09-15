package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestPerspectivePlausibleActions(t *testing.T) {
	remainingWorlds := BuildPossibleWorlds([]string{"abc", "aab"})
	perspective := &Perspective{}
	perspective.AvailableActions = BuildAvailableActions(remainingWorlds)
	perspective.ActionFilter = TestActionFilter

	actual := perspective.PlausibleActions(remainingWorlds)
	test.AssertDeepEqual(actual, []Action{"hi"}, t)
}

// func TestExploreReactions(t *testing.T) {
// 	worlds := BuildPossibleWorlds([]string{"cab", "aab"})
// 	aa := BuildAvailableActions(worlds)
// 	actual := ExploreReactions(worlds, aa.AllActions(), TestAgreement)
// 	test.AssertDeepEqual(actual.Get("a", "a"), []World{"aab", "cab"}, t)
// 	test.AssertDeepEqual(actual.Get("c", "c"), []World{"cab"}, t)
// 	test.AssertDeepEqual(actual.Get("c", "!"), []World{"aab"}, t)
// }
