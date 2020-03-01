package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestReactionExplorerGet(t *testing.T) {
	re := NewReactionExplorer()
	action := Action("a")
	response := Response("a--")
	// reaction := Reaction{Action: "a", Response: "a--", Reward: 0.75}
	actual := re.Get(action, response)
	test.AssertDeepEqual(actual, make([]World, 0), t)

	re.Add(action, response, 0.75, "abc")
	actual = re.Get(action, response)
	test.AssertDeepEqual(actual, []World{"abc"}, t)

	re.Add(action, response, 0.5, "aaa")
	actual = re.Get(action, response)
	test.AssertDeepEqual(actual, []World{"abc", "aaa"}, t)
}

func TestReactionExplorerActionCounts(t *testing.T) {
	re := NewReactionExplorer()
	action := Action("a")

	actual := re.ActionCounts(action)
	test.AssertDeepEqual(actual, []float64{}, t)

	re.Add(action, Response("a--"), 0.75, "abc")
	re.Add(action, Response("a--"), 1.0, "aaa")
	re.Add(action, Response("!"), 0.0, "bbc")

	actual = re.ActionCounts(action)
	test.AssertDeepEqual(actual, []float64{2.0, 1.0}, t)
}

func TestReactionExplorerActionTotal(t *testing.T) {
	re := NewReactionExplorer()
	action := Action("a")

	actual := re.ActionTotal(action)
	test.AssertDeepEqual(actual, 0.0, t)

	re.Add(action, Response("a--"), 0.75, "abc")
	re.Add(action, Response("a--"), 1.0, "aaa")
	re.Add(action, Response("!"), 0.0, "bbc")

	actual = re.ActionTotal(action)
	test.AssertDeepEqual(actual, 3.0, t)
}

func TestReactionExplorerActionRewardTotal(t *testing.T) {
	re := NewReactionExplorer()
	action := Action("a")

	actual := re.ActionRewardTotal(action)
	test.AssertDeepEqual(actual, 0.0, t)

	re.Add(action, Response("a--"), 0.75, "abc")
	re.Add(action, Response("a--"), 1.0, "aaa")
	re.Add(action, Response("!"), 0.2, "bbc")

	actual = re.ActionRewardTotal(action)
	test.AssertDeepEqual(actual, 1.95, t)
}
