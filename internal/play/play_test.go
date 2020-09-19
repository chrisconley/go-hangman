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

func TestDecideDeterminism(t *testing.T) {
	predictions := game.Predictions{
		"a": game.Prediction{Action: "a", Reward: 1, Entropy: 0},
		"b": game.Prediction{Action: "b", Reward: 1, Entropy: 0},
		"c": game.Prediction{Action: "c", Reward: 0, Entropy: 0},
		"d": game.Prediction{Action: "d", Reward: 0.5, Entropy: 1},
		"e": game.Prediction{Action: "e", Reward: 1, Entropy: 0},
		"f": game.Prediction{Action: "f", Reward: 0, Entropy: 0},
		"g": game.Prediction{Action: "g", Reward: 0, Entropy: 0},
		"h": game.Prediction{Action: "h", Reward: 0, Entropy: 0},
		"i": game.Prediction{Action: "i", Reward: 1, Entropy: 0},
		"j": game.Prediction{Action: "j", Reward: 0, Entropy: 0},
		"k": game.Prediction{Action: "k", Reward: 0, Entropy: 0},
		"l": game.Prediction{Action: "l", Reward: 0, Entropy: 0},
		"m": game.Prediction{Action: "m", Reward: 0, Entropy: 0},
		"n": game.Prediction{Action: "n", Reward: 0, Entropy: 0},
		"o": game.Prediction{Action: "o", Reward: 1, Entropy: 0},
		"p": game.Prediction{Action: "p", Reward: 0, Entropy: 0},
		"q": game.Prediction{Action: "q", Reward: 0, Entropy: 0},
		"r": game.Prediction{Action: "r", Reward: 0, Entropy: 0},
		"s": game.Prediction{Action: "s", Reward: 0.5, Entropy: 1},
		"t": game.Prediction{Action: "t", Reward: 0.5, Entropy: 1},
		"u": game.Prediction{Action: "u", Reward: 1, Entropy: 0},
		"v": game.Prediction{Action: "v", Reward: 0, Entropy: 0},
		"w": game.Prediction{Action: "w", Reward: 0, Entropy: 0},
		"x": game.Prediction{Action: "x", Reward: 0.5, Entropy: 1},
		"y": game.Prediction{Action: "y", Reward: 0, Entropy: 0},
		"z": game.Prediction{Action: "z", Reward: 0, Entropy: 0},
	}

	motivation := game.Motivation{Reward: 0.5, Entropy: 0.25, Time: 0.25}
	for i := 0; i < 100; i++ {
		actual := Decide(predictions, motivation)
		test.AssertDeepEqual(actual, game.Action("d"), t)
	}
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
