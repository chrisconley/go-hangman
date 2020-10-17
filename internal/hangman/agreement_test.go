package hangman

import (
	"fmt"
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/test"
)

func TestGetActions(t *testing.T) {
	game.BUniverse = game.InitUniverse([]string{"aa", "bb", "cc"})
	aa := game.BUniverse.Actions

	// When there are the same or more remaining worlds than letters, only use letters
	remainingWorlds := game.NewPossibleWorlds(2)
	lines := []string{"aa", "bb", "cc"}
	for index := range lines {
		remainingWorlds.Add(uint64(index))
	}
	gp := game.InitGamePlay()
	actions := GetActions(gp, remainingWorlds, aa)
	test.AssertDeepEqual(actions, game.BuildActions("a", "b", "c"), t)

	// When there are less remaining worlds than letters, add the worlds
	remainingWorlds = game.NewPossibleWorlds(1)
	lines = []string{"aa", "bb"}
	for index := range lines {
		remainingWorlds.Add(uint64(index))
	}
	gp = game.InitGamePlay()
	actions = GetActions(gp, remainingWorlds, aa)
	test.AssertDeepEqual(actions, game.BuildActions("a", "b", "c", "aa", "bb"), t)

	// When an action has already been used, don't add it
	remainingWorlds = game.NewPossibleWorlds(1)
	lines = []string{"aa", "bb", "cc"}
	for index := range lines {
		remainingWorlds.Add(uint64(index))
	}
	gp = game.InitGamePlay()
	gp.AddTurn(game.Reaction{Action: game.NewAction("b")})
	actions = GetActions(gp, remainingWorlds, aa)
	test.AssertDeepEqual(actions, game.BuildActions("a", "c"), t)

	// When a world action has already been used, don't add it
	remainingWorlds = game.NewPossibleWorlds(1)
	lines = []string{"aa", "bb"}
	for index := range lines {
		remainingWorlds.Add(uint64(index))
	}
	gp = game.InitGamePlay()
	gp.AddTurn(game.Reaction{Action: game.NewAction("bb")})
	actions = GetActions(gp, remainingWorlds, aa)
	test.AssertDeepEqual(actions, game.BuildActions("a", "b", "c", "aa"), t)
}

func TestAct(t *testing.T) {
	reaction := Act(game.InitGamePlay(), game.NewWorld(0, "abc"), game.NewAction("a"))
	test.AssertDeepEqual(reaction.Response, game.Response("a--"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = Act(game.InitGamePlay(), game.NewWorld(0, "abc"), game.NewAction("b"))
	test.AssertDeepEqual(reaction.Response, game.Response("-b-"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = Act(game.InitGamePlay(), game.NewWorld(0, "abc"), game.NewAction("c"))
	test.AssertDeepEqual(reaction.Response, game.Response("--c"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = Act(game.InitGamePlay(), game.NewWorld(0, "aaebaacdaa"), game.NewAction("a"))
	test.AssertDeepEqual(reaction.Response, game.Response("aa--aa--aa"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = Act(game.InitGamePlay(), game.NewWorld(0, "abc"), game.NewAction("abc"))
	test.AssertDeepEqual(reaction.Response, game.Response("abc"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = Act(game.InitGamePlay(), game.NewWorld(0, "abc"), game.NewAction("x"))
	test.AssertDeepEqual(reaction.Response, game.Response("!"), t)
	test.AssertDeepEqual(reaction.Reward, 0.0, t)

	reaction = Act(game.InitGamePlay(), game.NewWorld(0, "abc"), game.NewAction("aaa"))
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
		Act(game.InitGamePlay(), game.NewWorld(0, world), game.NewAction("x"))
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
		Act(game.InitGamePlay(), game.NewWorld(0, world), game.NewAction("x"))
	}
}
