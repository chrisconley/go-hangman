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
