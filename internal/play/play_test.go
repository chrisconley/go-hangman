package play

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"

	// "github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/test"
)

func TestDecide(t *testing.T) {
	predictions := game.Predictions{
		"a": game.Prediction{Action: "a", Entropy: 1.75, Reward: 0.9},
		"b": game.Prediction{Action: "b", Entropy: 2.50, Reward: 0.1},
	}

	// All Entropy
	motivation := game.Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
	actual := Decide(predictions, motivation)
	test.AssertDeepEqual(actual, game.Action("b"), t)

	// Weighted Reward with Entropy
	motivation = game.Motivation{Reward: 0.75, Entropy: 0.25, Time: 0.0}
	actual = Decide(predictions, motivation)
	test.AssertDeepEqual(actual, game.Action("a"), t)
}

func TestDecideLexicallySortsOnTie(t *testing.T) {
	predictions := game.Predictions{
		"b": game.Prediction{Action: "b", Entropy: 1.0, Reward: 0.5},
		"a": game.Prediction{Action: "a", Entropy: 1.0, Reward: 0.5},
	}

	// All Entropy
	motivation := game.Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
	actual := Decide(predictions, motivation)
	test.AssertDeepEqual(actual, game.Action("a"), t)
}

func TestAct(t *testing.T) {
	truth := game.World("abc")
	actual := Act(game.Action("a"), game.TestAgreement, truth)
	test.AssertDeepEqual(actual, game.Reaction{Action: "a", Response: "a", Reward: 1.0}, t)
}

func TestObserve(t *testing.T) {
	t.Skip("TODO: Fix test setup")
	perspective := &game.Perspective{}
	perspective.Knowledge = game.Knowledge{}

	reaction := game.Reaction{Response: "a", Reward: 1.0}
	Observe(game.Action("a"), reaction, perspective)
	test.AssertDeepEqual(perspective.Knowledge, game.Knowledge{reaction}, t)
}

func TestPlayBHangman(t *testing.T) {
	// lines := []string{
	// 	"abc",
	// 	"bba",
	// 	"cab",
	// }
	lines := []string{
		"aardvarks",
		"aasvogels",
		"abamperes",
		"abandoned",
		"abandoner",
		"abasement",
		"abashment",
		"abatement",
		"abattises",
		"abattoirs",
	}
	universeLength := uint64(len(lines))
	game.BUniverse = make([]game.BWorld, universeLength)
	pw := game.NewPossibleWorlds(universeLength)
	for index, name := range lines {
		game.BUniverse[index] = game.BWorld{Index: uint64(index), Name: name}
		pw.Bits.SetBit(uint64(index))
	}

	for _, line := range lines {
		worlds := game.BuildPossibleWorlds(lines)
		perspective := &game.Perspective{}
		perspective.Knowledge = game.Knowledge{}
		// perspective.PossibleWorlds = worlds
		perspective.BPossibleWorlds = pw
		perspective.AvailableActions = game.BuildAvailableActions(worlds)
		perspective.Agreement = hangman.HangmanAgreement
		perspective.ActionFilter = hangman.ActionFilter
		perspective.Motivation = game.Motivation{Reward: 0.5, Entropy: 0.25, Time: 0.25}
		truth := game.World(line)
		actual := PlayGame(truth, perspective, &game.NoMemory{})
		test.AssertDeepEqual(actual, truth, t)
	}
}

// func TestPlayHangman(t *testing.T) {
// 	worlds := game.BuildPossibleWorlds([]string{
// 		"abc",
// 		"aab",
// 		"bbb",
// 	})
// 	perspective := &game.Perspective{}
// 	perspective.Knowledge = game.Knowledge{}
// 	perspective.PossibleWorlds = worlds
// 	perspective.AvailableActions = game.BuildAvailableActions(worlds)
// 	perspective.Agreement = hangman.HangmanAgreement
// 	perspective.ActionFilter = hangman.ActionFilter
// 	perspective.Motivation = game.Motivation{0.5, 0.25, 0.25}

// 	truth := game.World("bbb")
// 	actual := PlayGame(truth, perspective, &game.NoMemory{})
// 	test.AssertDeepEqual(actual, truth, t)
// }

// func TestPlayMastermind(t *testing.T) {
// 	worlds := game.BuildPossibleWorlds([]string{
// 		"111",
// 		"112",
// 		"121",
// 		"211",
// 	})
// 	// perspective := NewPerspective(
// 	// 	worlds: worlds,
// 	// 	motivation: game.Motivation{0.5, 0.25, 0.25},
// 	// 	game: "mastermind"
// 	// )
// 	perspective := &game.Perspective{}
// 	perspective.Knowledge = game.Knowledge{}
// 	perspective.PossibleWorlds = worlds
// 	perspective.AvailableActions = game.BuildAvailableActions(worlds)
// 	perspective.Agreement = mastermind.Agreement
// 	perspective.ActionFilter = mastermind.ActionFilter
// 	perspective.Motivation = game.Motivation{0.5, 0.25, 0.25}

// 	truth := game.World("111")
// 	actual := PlayGame(truth, perspective, &game.NoMemory{})
// 	test.AssertDeepEqual(actual, truth, t)
// }
