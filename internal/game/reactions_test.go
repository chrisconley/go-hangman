package game

import (
	"fmt"
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestReactionExplorerGet(t *testing.T) {
	re := NewReactionExplorer()
	action := Action("a")
	response := Response("a--")
	actual := re.Get(action, response)
	test.AssertDeepEqual(actual.Worlds, make([]World, 0), t)

	re.Add(action, response, 0.75, "abc")
	actual = re.Get(action, response)
	test.AssertDeepEqual(actual.Worlds, []World{"abc"}, t)

	re.Add(action, response, 0.5, "aaa")
	actual = re.Get(action, response)
	test.AssertDeepEqual(actual.Worlds, []World{"abc", "aaa"}, t)
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

func TestReactionExplorerActions(t *testing.T) {
	re := NewReactionExplorer()
	test.AssertDeepEqual(re.Actions(), []Action{}, t)

	re.Add(Action("a"), Response("a--"), 0.75, "abc")
	re.Add(Action("a"), Response("a--"), 1.0, "aaa")
	re.Add(Action("b"), Response("!"), 0.0, "aac")

	test.AssertDeepEqual(re.Actions(), []Action{"a", "b"}, t)
}

func BenchmarkActionsSmall(b *testing.B) {
	re := NewReactionExplorer()
	numElements := 10
	for i := 0; i < numElements; i++ {
		re.Add(Action(fmt.Sprintf("%d", i)), "!", 0.0, "a")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		re.Actions()
	}
}

func BenchmarkActionsLarge(b *testing.B) {
	re := NewReactionExplorer()
	numElements := 10000
	for i := 0; i < numElements; i++ {
		re.Add(Action(fmt.Sprintf("%d", i)), "!", 0.0, "a")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		re.Actions()
	}
}

// func BenchmarkAllLarge(b *testing.B) {
// 	gi := NewGroupedIncrementer()
// 	numElements := 1000000
// 	for i := 0; i < numElements; i++ {
// 		gi.Increment(fmt.Sprintf("%d", i), 2.4)
// 	}
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		gi.All()
// 	}
// }
