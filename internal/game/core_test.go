package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestWorld(t *testing.T) {
	world := NewWorld("hi")
	test.AssertDeepEqual(world.String(), "hi", t)
	test.AssertDeepEqual(world.Length(), 2, t)
	test.AssertDeepEqual(world.Bytes(), []byte("hi"), t)
	test.AssertDeepEqual(world.Hash(), uint64(628919584683901914), t)
	test.AssertDeepEqual(world.Matches(NewWorld("hi")), true, t)
	test.AssertDeepEqual(world.Matches(NewWorld("h")), false, t)
}

func TestAction(t *testing.T) {
	action := NewAction("hi")
	test.AssertDeepEqual(action.String(), "hi", t)
	test.AssertDeepEqual(action.Length(), 2, t)
	test.AssertDeepEqual(action.Bytes(), []byte("hi"), t)
	test.AssertDeepEqual(action.Hash(), uint64(628919584683901914), t)
	test.AssertDeepEqual(action.Matches(NewAction("hi")), true, t)
	test.AssertDeepEqual(action.Matches(NewAction("h")), false, t)
}

func TestResponse(t *testing.T) {
	response := NewResponse("hi")
	test.AssertDeepEqual(response.String(), "hi", t)
	test.AssertDeepEqual(response.Length(), 2, t)
	test.AssertDeepEqual(response.Bytes(), []byte("hi"), t)
	test.AssertDeepEqual(response.Hash(), uint64(628919584683901914), t)
	test.AssertDeepEqual(response.Matches(NewResponse("hi")), true, t)
	test.AssertDeepEqual(response.Matches(NewResponse("h")), false, t)
}
