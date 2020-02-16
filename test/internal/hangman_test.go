package test

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
)

func TestUpdatePossibleActions(t *testing.T) {
	reaction := game.Reaction{
		Action: game.Action("a"),
		Response: game.Response("a--"),
		Reward: 1.0,
	}
	possibleActions := game.PossibleActions{"a", "b", "c"}

	actual := game.UpdatePossibleActions(reaction, possibleActions)
	expected := game.PossibleActions{"b", "c"}
	assertDeepEqual(actual, expected, t)
}

func TestUpdatePossibleWorldsPositiveResponse(t *testing.T) {
	reaction := game.Reaction{
		Action: game.Action("a"),
		Response: game.Response("a--"),
		Reward: 1.0,
	}
	possibleWorlds := game.PossibleWorlds{
		game.PossibleWorld{"abc", 0.5},
		game.PossibleWorld{"aab", 0.5},
		game.PossibleWorld{"bab", 0.5},
		game.PossibleWorld{"bbb", 0.5},
	}

	actual := game.UpdatePossibleWorlds(reaction, possibleWorlds)
	expected := game.PossibleWorlds{
		game.PossibleWorld{"abc", 0.5},
		game.PossibleWorld{"aab", 0.5},
	}
	assertDeepEqual(actual, expected, t)
}

func TestUpdatePossibleWorldsNegativeResponse(t *testing.T) {
	reaction := game.Reaction{
		Action: game.Action("a"),
		Response: game.Response("!"),
		Reward: 1.0,
	}
	possibleWorlds := game.PossibleWorlds{
		game.PossibleWorld{"abc", 0.5},
		game.PossibleWorld{"aab", 0.5},
		game.PossibleWorld{"bab", 0.5},
		game.PossibleWorld{"bbb", 0.5},
	}

	actual := game.UpdatePossibleWorlds(reaction, possibleWorlds)
	expected := game.PossibleWorlds{
		game.PossibleWorld{"bbb", 0.5},
	}
	assertDeepEqual(actual, expected, t)
}

func TestHangmanAgreement(t *testing.T) {
	reaction := game.HangmanAgreement(game.World("abc"), game.Action("a"))
	assertDeepEqual(reaction.Response, game.Response("a--"), t)
	assertDeepEqual(reaction.Reward, 1.0, t)

	reaction = game.HangmanAgreement(game.World("abc"), game.Action("b"))
	assertDeepEqual(reaction.Response, game.Response("-b-"), t)
	assertDeepEqual(reaction.Reward, 1.0, t)

	reaction = game.HangmanAgreement(game.World("abc"), game.Action("c"))
	assertDeepEqual(reaction.Response, game.Response("--c"), t)
	assertDeepEqual(reaction.Reward, 1.0, t)

	reaction = game.HangmanAgreement(game.World("aaebaacdaa"), game.Action("a"))
	assertDeepEqual(reaction.Response, game.Response("aa--aa--aa"), t)
	assertDeepEqual(reaction.Reward, 1.0, t)

	reaction = game.HangmanAgreement(game.World("abc"), game.Action("abc"))
	assertDeepEqual(reaction.Response, game.Response("abc"), t)
	assertDeepEqual(reaction.Reward, 1.0, t)

	reaction = game.HangmanAgreement(game.World("abc"), game.Action("x"))
	assertDeepEqual(reaction.Response, game.Response("!"), t)
	assertDeepEqual(reaction.Reward, 0.0, t)

	reaction = game.HangmanAgreement(game.World("abc"), game.Action("aaa"))
	assertDeepEqual(reaction.Response, game.Response("!"), t)
	assertDeepEqual(reaction.Reward, 0.0, t)
}
