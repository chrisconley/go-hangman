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

func TestActRewards(t *testing.T) {
	reaction := Act(game.World("1122"), game.Action("1313"))
	test.AssertDeepEqual(reaction.Response, game.Response("BW"), t)
	test.AssertAlmostEqual(reaction.Reward, 0.375, 4, t)

	reaction = Act(game.World("1122"), game.Action("2211"))
	test.AssertDeepEqual(reaction.Response, game.Response("WWWW"), t)
	test.AssertAlmostEqual(reaction.Reward, 0.5, 4, t)

	reaction = Act(game.World("11224"), game.Action("13135"))
	test.AssertDeepEqual(reaction.Response, game.Response("BW"), t)
	test.AssertAlmostEqual(reaction.Reward, 0.3, 4, t)

	reaction = Act(game.World("112244"), game.Action("131355"))
	test.AssertDeepEqual(reaction.Response, game.Response("BW"), t)
	test.AssertAlmostEqual(reaction.Reward, 0.25, 4, t)
}

func TestActWikipediaExample(t *testing.T) {
	// Wikipedia example: https://en.wikipedia.org/wiki/Mastermind_(board_game)#Gameplay_and_rules
	reaction := Act(game.World("1122"), game.Action("1112"))
	test.AssertDeepEqual(reaction.Response, game.Response("BBB"), t)
	test.AssertDeepEqual(reaction.Reward, 0.75, t)
}

func TestActWhiteCounting(t *testing.T) {
	// Duplicates in world
	reaction := Act(game.World("1112"), game.Action("7891"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	reaction = Act(game.World("2111"), game.Action("1789"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	reaction = Act(game.World("1211"), game.Action("7189"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	// Duplicates in action
	reaction = Act(game.World("1234"), game.Action("7111"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	reaction = Act(game.World("2341"), game.Action("1117"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	reaction = Act(game.World("2134"), game.Action("1711"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	// Duplicates in both
	reaction = Act(game.World("2311"), game.Action("1178"))
	test.AssertDeepEqual(reaction.Response, game.Response("WW"), t)

	reaction = Act(game.World("1123"), game.Action("7811"))
	test.AssertDeepEqual(reaction.Response, game.Response("WW"), t)

	reaction = Act(game.World("2113"), game.Action("1781"))
	test.AssertDeepEqual(reaction.Response, game.Response("WW"), t)

	reaction = Act(game.World("1231"), game.Action("7118"))
	test.AssertDeepEqual(reaction.Response, game.Response("WW"), t)
}

func TestActInstructionsExamples(t *testing.T) {
	// Instructions examples from docs/mastermindrules.pdf
	// https://magisterrex.files.wordpress.com/2014/07/mastermindrules.pdf

	// Basic example responding with "BW"
	// "B" is for the third position "U" guess
	// "W" is for the first position "W" guess
	reaction := Act(game.World("RWUG"), game.Action("WYUK"))
	test.AssertDeepEqual(reaction.Response, game.Response("BW"), t)

	// Prefer "B" response over "W".
	// Do not respond with "BW" because the code word only has one G
	reaction = Act(game.World("RUGK"), game.Action("WYGG"))
	test.AssertDeepEqual(reaction.Response, game.Response("B"), t)

	// Only one "W" because the guess only has one "Y"
	reaction = Act(game.World("YUYK"), game.Action("RYGW"))
	test.AssertDeepEqual(reaction.Response, game.Response("W"), t)

	// Two "R" guesses; one is in right place, one is not
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
