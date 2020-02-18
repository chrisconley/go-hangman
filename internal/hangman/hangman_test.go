package hangman

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/test"
)

func TestUpdatePossibleActions(t *testing.T) {
	reaction := game.Reaction{
		Action:   game.Action("a"),
		Response: game.Response("a--"),
		Reward:   1.0,
	}
	possibleActions := game.PossibleActions{"a", "b", "c"}

	actual := UpdatePossibleActions(reaction, possibleActions)
	expected := game.PossibleActions{"b", "c"}
	test.AssertDeepEqual(actual, expected, t)
}

func TestUpdatePossibleWorldsPositiveResponse(t *testing.T) {
	reaction := game.Reaction{
		Action:   game.Action("a"),
		Response: game.Response("a--"),
		Reward:   1.0,
	}
	possibleWorlds := game.PossibleWorlds{
		game.PossibleWorld{"abc", 0.5},
		game.PossibleWorld{"aab", 0.5},
		game.PossibleWorld{"bab", 0.5},
		game.PossibleWorld{"bbb", 0.5},
	}

	actual := UpdatePossibleWorlds(reaction, possibleWorlds)
	expected := game.PossibleWorlds{
		game.PossibleWorld{"abc", 0.5},
		game.PossibleWorld{"aab", 0.5},
	}
	test.AssertDeepEqual(actual, expected, t)
}

func TestUpdatePossibleWorldsNegativeResponse(t *testing.T) {
	reaction := game.Reaction{
		Action:   game.Action("a"),
		Response: game.Response("!"),
		Reward:   1.0,
	}
	possibleWorlds := game.PossibleWorlds{
		game.PossibleWorld{"abc", 0.5},
		game.PossibleWorld{"aab", 0.5},
		game.PossibleWorld{"bab", 0.5},
		game.PossibleWorld{"bbb", 0.5},
	}

	actual := UpdatePossibleWorlds(reaction, possibleWorlds)
	expected := game.PossibleWorlds{
		game.PossibleWorld{"bbb", 0.5},
	}
	test.AssertDeepEqual(actual, expected, t)
}

func TestHangmanAgreement(t *testing.T) {
	reaction := HangmanAgreement(game.World("abc"), game.Action("a"))
	test.AssertDeepEqual(reaction.Response, game.Response("a--"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = HangmanAgreement(game.World("abc"), game.Action("b"))
	test.AssertDeepEqual(reaction.Response, game.Response("-b-"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = HangmanAgreement(game.World("abc"), game.Action("c"))
	test.AssertDeepEqual(reaction.Response, game.Response("--c"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = HangmanAgreement(game.World("aaebaacdaa"), game.Action("a"))
	test.AssertDeepEqual(reaction.Response, game.Response("aa--aa--aa"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = HangmanAgreement(game.World("abc"), game.Action("abc"))
	test.AssertDeepEqual(reaction.Response, game.Response("abc"), t)
	test.AssertDeepEqual(reaction.Reward, 1.0, t)

	reaction = HangmanAgreement(game.World("abc"), game.Action("x"))
	test.AssertDeepEqual(reaction.Response, game.Response("!"), t)
	test.AssertDeepEqual(reaction.Reward, 0.0, t)

	reaction = HangmanAgreement(game.World("abc"), game.Action("aaa"))
	test.AssertDeepEqual(reaction.Response, game.Response("!"), t)
	test.AssertDeepEqual(reaction.Reward, 0.0, t)
}
