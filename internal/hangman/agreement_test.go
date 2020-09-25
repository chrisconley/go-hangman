package hangman

import (
	"fmt"
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/test"
)

func TestGetActions(t *testing.T) {
	allWorlds := game.BuildPossibleWorlds([]string{"aa", "bb"})
	aa := game.BuildAvailableActions(allWorlds)

	remainingWorlds := game.BuildPossibleWorlds([]string{"aa", "bb"})
	actions := GetActions(remainingWorlds, aa)
	test.AssertDeepEqual(actions, []game.Action{"a", "b"}, t)

	remainingWorlds = game.BuildPossibleWorlds([]string{"aa"})
	actions = GetActions(remainingWorlds, aa)
	test.AssertDeepEqual(actions, []game.Action{"a", "b", "aa"}, t)
}

func TestAct(t *testing.T) {
	reaction := Act(game.World("abc"), game.Action("a"))
	test.AssertDeepEqual(reaction.Response, game.Response("a--"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = Act(game.World("abc"), game.Action("b"))
	test.AssertDeepEqual(reaction.Response, game.Response("-b-"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = Act(game.World("abc"), game.Action("c"))
	test.AssertDeepEqual(reaction.Response, game.Response("--c"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = Act(game.World("aaebaacdaa"), game.Action("a"))
	test.AssertDeepEqual(reaction.Response, game.Response("aa--aa--aa"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = Act(game.World("abc"), game.Action("abc"))
	test.AssertDeepEqual(reaction.Response, game.Response("abc"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = Act(game.World("abc"), game.Action("x"))
	test.AssertDeepEqual(reaction.Response, game.Response("!"), t)
	test.AssertDeepEqual(reaction.Reward, 0.0, t)

	reaction = Act(game.World("abc"), game.Action("aaa"))
	test.AssertDeepEqual(reaction.Response, game.Response("!"), t)
	test.AssertDeepEqual(reaction.Reward, 0.0, t)
}

func BenchmarkActSmall(b *testing.B) {
	world := ""
	numElements := 10
	for i := 0; i < numElements; i++ {
		world = world + fmt.Sprintf("hi%d", i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Act(game.World(world), game.Action("x"))
	}
}

func BenchmarkActLarge(b *testing.B) {
	world := ""
	numElements := 1000
	for i := 0; i < numElements; i++ {
		world = world + fmt.Sprintf("hi%d", i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Act(game.World(world), game.Action("x"))
	}
}
