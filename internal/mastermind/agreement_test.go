package mastermind

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/test"
)

func TestGetActions(t *testing.T) {
	allWorlds := game.BuildPossibleWorlds([]string{"aa", "bb"})
	aa := game.BuildAvailableActions(allWorlds)

	remainingWorlds := game.BuildPossibleWorlds([]string{"aa", "bb"})
	actions := GetActions(remainingWorlds, aa)
	test.AssertDeepEqual(actions, []game.Action{"aa", "bb"}, t)

	remainingWorlds = game.BuildPossibleWorlds([]string{"aa"})
	actions = GetActions(remainingWorlds, aa)
	test.AssertDeepEqual(actions, []game.Action{"aa", "bb"}, t)
}

func TestAct(t *testing.T) {
	reaction := Act(game.World("1122"), game.Action("1313"))
	test.AssertDeepEqual(reaction.Response, game.Response("BW"), t)
	test.AssertDeepEqual(reaction.Reward, 0.375, t)

	reaction = Act(game.World("1122"), game.Action("2211"))
	test.AssertDeepEqual(reaction.Response, game.Response("WWWW"), t)
	test.AssertDeepEqual(reaction.Reward, 0.5, t)
}

func TestActWikipediaExample(t *testing.T) {
	// Wikipedia example: https://en.wikipedia.org/wiki/Mastermind_(board_game)#Gameplay_and_rules
	reaction := Act(game.World("1122"), game.Action("1112"))
	test.AssertDeepEqual(reaction.Response, game.Response("BBB"), t)
	test.AssertDeepEqual(reaction.Reward, 0.75, t)
}

func TestActWhiteCounting(t *testing.T) {
	// Counting whites
	reaction := Act(game.World("1112"), game.Action("3331"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	reaction = Act(game.World("2111"), game.Action("1333"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	reaction = Act(game.World("1222"), game.Action("3111"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	reaction = Act(game.World("2221"), game.Action("1113"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)
}

func TestActInstructionsExamples(t *testing.T) {
	// Instructions examples fromj
	// https://magisterrex.files.wordpress.com/2014/07/mastermindrules.pdf
	reaction := Act(game.World("RWUG"), game.Action("WYUK"))
	test.AssertDeepEqual(reaction.Response, game.Response("BW"), t)

	reaction = Act(game.World("RUGK"), game.Action("WYGG"))
	test.AssertDeepEqual(reaction.Response, game.Response("B"), t)

	reaction = Act(game.World("YUYK"), game.Action("RYGW"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	reaction = Act(game.World("RURY"), game.Action("GRRW"))
	test.AssertDeepEqual(reaction.Response, game.Response("BW"), t)
}

func TestActFromPythonImplementation(t *testing.T) {
	// Various examples using 3632 code word from python implementation tests
	reaction := Act(game.World("3632"), game.Action("1611"))
	test.AssertDeepEqual(reaction.Response, game.Response("B"), t)
	test.AssertDeepEqual(reaction.Reward, 0.25, t)

	reaction = Act(game.World("3632"), game.Action("6642"))
	test.AssertDeepEqual(reaction.Response, game.Response("BB"), t)

	reaction = Act(game.World("3632"), game.Action("6326"))
	test.AssertDeepEqual(reaction.Response, game.Response("WWW"), t)

	reaction = Act(game.World("3632"), game.Action("2363"))
	test.AssertDeepEqual(reaction.Response, game.Response("WWWW"), t)

	reaction = Act(game.World("3632"), game.Action("6323"))
	test.AssertDeepEqual(reaction.Response, game.Response("WWWW"), t)

	reaction = Act(game.World("3632"), game.Action("1435"))
	test.AssertDeepEqual(reaction.Response, game.Response("B"), t)

	reaction = Act(game.World("3632"), game.Action("1335"))
	test.AssertDeepEqual(reaction.Response, game.Response("BW"), t)
}
