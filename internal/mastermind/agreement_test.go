package mastermind

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/test"
)
func TestAgreement(t *testing.T) {
	reaction := Agreement(game.World("1122"), game.Action("1313"))
	test.AssertDeepEqual(reaction.Response, game.Response("BW"), t)
	test.AssertDeepEqual(reaction.Reward, 0.375, t)

	reaction = Agreement(game.World("1122"), game.Action("2211"))
	test.AssertDeepEqual(reaction.Response, game.Response("WWWW"), t)
	test.AssertDeepEqual(reaction.Reward, 0.5, t)
}

func TestAgreementWikipediaExample(t *testing.T) {
	// Wikipedia example: https://en.wikipedia.org/wiki/Mastermind_(board_game)#Gameplay_and_rules
	reaction := Agreement(game.World("1122"), game.Action("1112"))
	test.AssertDeepEqual(reaction.Response, game.Response("BBB"), t)
	test.AssertDeepEqual(reaction.Reward, 0.75, t)
}

func TestAgreementWhiteCounting(t *testing.T) {
	// Counting whites
	reaction := Agreement(game.World("1112"), game.Action("3331"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	reaction = Agreement(game.World("2111"), game.Action("1333"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	reaction = Agreement(game.World("1222"), game.Action("3111"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	reaction = Agreement(game.World("2221"), game.Action("1113"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)
}

func TestAgreementInstructionsExamples(t *testing.T) {
	// Instructions examples fromj
	// https://magisterrex.files.wordpress.com/2014/07/mastermindrules.pdf
	reaction := Agreement(game.World("RWUG"), game.Action("WYUK"))
	test.AssertDeepEqual(reaction.Response, game.Response("BW"), t)

	reaction = Agreement(game.World("RUGK"), game.Action("WYGG"))
	test.AssertDeepEqual(reaction.Response, game.Response("B"), t)

	reaction = Agreement(game.World("YUYK"), game.Action("RYGW"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	reaction = Agreement(game.World("RURY"), game.Action("GRRW"))
	test.AssertDeepEqual(reaction.Response, game.Response("BW"), t)
}

func TestAgreementFromPythonImplementation(t *testing.T) {
	// Various examples using 3632 code word from python implementation tests
	reaction := Agreement(game.World("3632"), game.Action("1611"))
	test.AssertDeepEqual(reaction.Response, game.Response("B"), t)
	test.AssertDeepEqual(reaction.Reward, 0.25, t)

	reaction = Agreement(game.World("3632"), game.Action("6642"))
	test.AssertDeepEqual(reaction.Response, game.Response("BB"), t)

	reaction = Agreement(game.World("3632"), game.Action("6326"))
	test.AssertDeepEqual(reaction.Response, game.Response("WWW"), t)

	reaction = Agreement(game.World("3632"), game.Action("2363"))
	test.AssertDeepEqual(reaction.Response, game.Response("WWWW"), t)

	reaction = Agreement(game.World("3632"), game.Action("6323"))
	test.AssertDeepEqual(reaction.Response, game.Response("WWWW"), t)

	reaction = Agreement(game.World("3632"), game.Action("1435"))
	test.AssertDeepEqual(reaction.Response, game.Response("B"), t)

	reaction = Agreement(game.World("3632"), game.Action("1335"))
	test.AssertDeepEqual(reaction.Response, game.Response("BW"), t)
}