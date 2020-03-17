package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestGameState(t *testing.T) {
	gs := GameState{}
	actual := gs.Key()
	test.AssertDeepEqual(actual, "", t)

	gs.Add(Reaction{Action: "a", Response: "a-"})
	actual = gs.Key()
	test.AssertDeepEqual(actual, "|a:a-", t)

	gs.Add(Reaction{Action: "b", Response: "!"})
	actual = gs.Key()
	test.AssertDeepEqual(actual, "|a:a-|b:!", t)

	gs.Add(Reaction{Action: "c", Response: "-c"})
	actual = gs.Key()
	test.AssertDeepEqual(actual, "|a:a-|b:!|c:-c", t)
}

func TestNoMemory(t *testing.T) {
	noMemory := &NoMemory{}
	gs := GameState{}
	gs.Add(Reaction{Action: "a", Response: "a-"})

	actualValue, actualResult := noMemory.Get(gs)
	test.AssertDeepEqual(actualValue, []World{}, t)
	test.AssertDeepEqual(actualResult, false, t)

	noMemory.Add(gs, []World{"abc"})
	actualValue, actualResult = noMemory.Get(gs)
	test.AssertDeepEqual(actualValue, []World{}, t)
	test.AssertDeepEqual(actualResult, false, t)
}

func TestNaiveMemory(t *testing.T) {
	naiveMemory := InitMemory(&NaiveMemory{})
	gs := GameState{}

	actualValue, actualResult := naiveMemory.Get(gs)
	test.AssertDeepEqual(actualValue, []World{}, t)
	test.AssertDeepEqual(actualResult, false, t)

	gs.Add(Reaction{Action: "a", Response: "a-"})
	actualValue, actualResult = naiveMemory.Get(gs)
	test.AssertDeepEqual(actualValue, []World{}, t)
	test.AssertDeepEqual(actualResult, false, t)

	naiveMemory.Add(gs, []World{"abc"})
	actualValue, actualResult = naiveMemory.Get(gs)
	test.AssertDeepEqual(actualValue, []World{"abc"}, t)
	test.AssertDeepEqual(actualResult, true, t)
}
