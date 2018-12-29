package test

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
)

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
