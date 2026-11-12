package mastermind

import (
	"sort"
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/chrisconley/go-hangman/test"
)

func TestGetActions(t *testing.T) {
	universe := game.InitUniverse([]string{"aa", "ab", "bb"})
	aa := universe.Actions

	remainingWorlds := game.NewPossibleWorlds(2)
	lines := []string{"aa", "ab", "bb"}
	for index := range lines {
		remainingWorlds.Add(uint64(index))
	}
	gp := game.InitGamePlay(0.0)
	actions := GetActions(gp, remainingWorlds, aa)
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].String() < actions[j].String()
	})
	test.AssertDeepEqual(actions, game.BuildActions("aa", "ab"), t)

	// Still return worlds that are not in remaining possible worlds
	remainingWorlds = game.NewPossibleWorlds(1)
	lines = []string{"aa"}
	for index := range lines {
		remainingWorlds.Add(uint64(index))
	}
	gp = game.InitGamePlay(0.0)
	actions = GetActions(gp, remainingWorlds, aa)
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].String() < actions[j].String()
	})
	test.AssertDeepEqual(actions, game.BuildActions("aa", "ab"), t)

	// Don't return worlds that have already been used
	remainingWorlds = game.NewPossibleWorlds(1)
	lines = []string{"aa"}
	for index := range lines {
		remainingWorlds.Add(uint64(index))
	}
	gp = game.InitGamePlay(0.0)
	gp.AddTurn(game.NewReaction(game.NewAction("aa"), "-", 0.0))
	actions = GetActions(gp, remainingWorlds, aa)
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].String() < actions[j].String()
	})
	test.AssertDeepEqual(actions, game.BuildActions("ab", "bb"), t)
}

func TestGetActionsUniqueness(t *testing.T) {
	lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")
	game.BUniverse = game.InitUniverse(lines)
	gp := game.InitGamePlay(game.BUniverse.PossibleWorlds.GetHash())
	actions := GetActions(gp, game.BUniverse.PossibleWorlds, game.BUniverse.Actions)
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].String() < actions[j].String()
	})
	expected := []game.Action{
		game.NewAction("AAAA"),
		game.NewAction("AAAB"),
		game.NewAction("AABB"),
		game.NewAction("AABC"),
		game.NewAction("ABCD"),
	}
	test.AssertDeepEqual(actions, expected, t)
}

func TestGetActionsUniquenessAfterFirstGuess(t *testing.T) {
	lines := utils.ReadLinesFromFilePath("../../fixtures/mastermind-ABCDEFG-4.txt")
	game.BUniverse = game.InitUniverse(lines)
	gp := game.InitGamePlay(game.BUniverse.PossibleWorlds.GetHash())
	gp.AddTurn(game.NewReaction(game.NewAction("ABCD"), "WWWW", 0.0))
	actions := GetActions(gp, game.BUniverse.PossibleWorlds, game.BUniverse.Actions)
	actionsMap := make(map[string]int8)
	for _, action := range actions {
		actionsMap[action.String()] = 1
	}
	_, ok := actionsMap["ABCE"]
	test.AssertDeepEqual(true, ok, t)
	_, ok = actionsMap["ABCF"]
	test.AssertDeepEqual(false, ok, t)
	_, ok = actionsMap["ABCG"]
	test.AssertDeepEqual(false, ok, t)

	_, ok = actionsMap["ABEE"]
	test.AssertDeepEqual(true, ok, t)
	_, ok = actionsMap["ABFF"]
	test.AssertDeepEqual(false, ok, t)
	_, ok = actionsMap["ABGG"]
	test.AssertDeepEqual(false, ok, t)

	_, ok = actionsMap["ABEF"]
	test.AssertDeepEqual(true, ok, t)
	_, ok = actionsMap["ABEG"]
	test.AssertDeepEqual(false, ok, t)
	_, ok = actionsMap["ABFE"]
	test.AssertDeepEqual(false, ok, t)
	_, ok = actionsMap["ABFG"]
	test.AssertDeepEqual(false, ok, t)
	_, ok = actionsMap["ABGE"]
	test.AssertDeepEqual(false, ok, t)
	_, ok = actionsMap["ABGF"]
	test.AssertDeepEqual(false, ok, t)
}

func TestActRewards(t *testing.T) {
	reaction := Act(game.InitGamePlay(0.0), game.NewWorld("1122"), game.NewAction("1313"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("BW"), t)
	test.AssertAlmostEqual(reaction.Reward, 0.375, 4, t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("1122"), game.NewAction("2211"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("WWWW"), t)
	test.AssertAlmostEqual(reaction.Reward, 0.5, 4, t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("11224"), game.NewAction("13135"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("BW"), t)
	test.AssertAlmostEqual(reaction.Reward, 0.3, 4, t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("112244"), game.NewAction("131355"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("BW"), t)
	test.AssertAlmostEqual(reaction.Reward, 0.25, 4, t)
}

func TestActMiss(t *testing.T) {
	reaction := Act(game.InitGamePlay(0.0), game.NewWorld("1122"), game.NewAction("3333"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("!"), t)
}

func TestActWikipediaExample(t *testing.T) {
	// Wikipedia example: https://en.wikipedia.org/wiki/Mastermind_(board_game)#Gameplay_and_rules
	reaction := Act(game.InitGamePlay(0.0), game.NewWorld("1122"), game.NewAction("1112"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("BBB"), t)
	test.AssertDeepEqual(reaction.Reward, 0.75, t)
}

func TestActWhiteCounting(t *testing.T) {
	// Duplicates in world
	reaction := Act(game.InitGamePlay(0.0), game.NewWorld("1112"), game.NewAction("7891"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("W"), t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("2111"), game.NewAction("1789"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("W"), t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("1211"), game.NewAction("7189"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("W"), t)

	// Duplicates in action
	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("1234"), game.NewAction("7111"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("W"), t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("2341"), game.NewAction("1117"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("W"), t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("2134"), game.NewAction("1711"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("W"), t)

	// Duplicates in both
	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("BCAA"), game.NewAction("AAYZ"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("WW"), t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("1123"), game.NewAction("7811"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("WW"), t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("2113"), game.NewAction("1781"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("WW"), t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("1231"), game.NewAction("7118"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("WW"), t)
}

func TestActInstructionsExamples(t *testing.T) {
	// Instructions examples from the Mastermind rules sheet (see REFERENCES.md)
	// https://magisterrex.files.wordpress.com/2014/07/mastermindrules.pdf

	// Basic example responding with "BW"
	// "B" is for the third position "U" guess
	// "W" is for the first position "W" guess
	reaction := Act(game.InitGamePlay(0.0), game.NewWorld("RWUG"), game.NewAction("WYUK"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("BW"), t)

	// Prefer "B" response over "W".
	// Do not respond with "BW" because the code word only has one G
	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("RUGK"), game.NewAction("WYGG"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("B"), t)

	// Only one "W" because the guess only has one "Y"
	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("YUYK"), game.NewAction("RYGW"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("W"), t)

	// Two "R" guesses; one is in right place, one is not
	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("RURY"), game.NewAction("GRRW"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("BW"), t)
}

func TestActFromPythonImplementation(t *testing.T) {
	// Various examples using 3632 code word from python implementation tests
	reaction := Act(game.InitGamePlay(0.0), game.NewWorld("3632"), game.NewAction("1611"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("B"), t)
	test.AssertDeepEqual(reaction.Reward, 0.25, t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("3632"), game.NewAction("6642"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("BB"), t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("3632"), game.NewAction("6326"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("WWW"), t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("3632"), game.NewAction("2363"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("WWWW"), t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("3632"), game.NewAction("6323"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("WWWW"), t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("3632"), game.NewAction("1435"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("B"), t)

	reaction = Act(game.InitGamePlay(0.0), game.NewWorld("3632"), game.NewAction("1335"))
	test.AssertDeepEqual(reaction.Response, game.NewResponse("BW"), t)
}

func TestKnuthCounts(t *testing.T) {
	// This verifies that this implementation follows the example in
	// Knuth's paper (see REFERENCES.md) with the same counts given in Figure 1.
	lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")
	groupedResponses := buildGroupedResponses(lines, "AABB")
	test.AssertDeepEqual(len(groupedResponses["WWWW"]), 1, t)
	test.AssertDeepEqual(len(groupedResponses["WWW"]), 16, t)
	test.AssertDeepEqual(len(groupedResponses["WW"]), 96, t)
	test.AssertDeepEqual(len(groupedResponses["W"]), 256, t)
	test.AssertDeepEqual(len(groupedResponses["!"]), 256, t)

	test.AssertDeepEqual(len(groupedResponses["BWWW"]), 0, t)
	test.AssertDeepEqual(len(groupedResponses["BWW"]), 36, t)
	test.AssertDeepEqual(len(groupedResponses["BW"]), 208, t)
	test.AssertDeepEqual(len(groupedResponses["B"]), 256, t)

	test.AssertDeepEqual(len(groupedResponses["BBWW"]), 4, t)
	test.AssertDeepEqual(len(groupedResponses["BBW"]), 32, t)
	test.AssertDeepEqual(len(groupedResponses["BB"]), 114, t)

	test.AssertDeepEqual(len(groupedResponses["BBBW"]), 0, t)
	test.AssertDeepEqual(len(groupedResponses["BBB"]), 20, t)

	test.AssertDeepEqual(len(groupedResponses["BBBB"]), 1, t)

	lines2 := groupedResponses["B"]
	groupedResponses2 := buildGroupedResponses(lines2, "ACDD")
	test.AssertDeepEqual(len(groupedResponses2["WWWW"]), 0, t)
	test.AssertDeepEqual(len(groupedResponses2["WWW"]), 7, t)
	test.AssertDeepEqual(len(groupedResponses2["WW"]), 41, t)
	test.AssertDeepEqual(len(groupedResponses2["W"]), 44, t)
	test.AssertDeepEqual(len(groupedResponses2["!"]), 16, t)

	test.AssertDeepEqual(len(groupedResponses2["BWWW"]), 2, t)
	test.AssertDeepEqual(len(groupedResponses2["BWW"]), 21, t)
	test.AssertDeepEqual(len(groupedResponses2["BW"]), 42, t)
	test.AssertDeepEqual(len(groupedResponses2["B"]), 34, t)

	test.AssertDeepEqual(len(groupedResponses2["BBWW"]), 3, t)
	test.AssertDeepEqual(len(groupedResponses2["BBW"]), 13, t)
	test.AssertDeepEqual(len(groupedResponses2["BB"]), 23, t)

	test.AssertDeepEqual(len(groupedResponses2["BBBW"]), 0, t)
	test.AssertDeepEqual(len(groupedResponses2["BBB"]), 9, t)

	test.AssertDeepEqual(len(groupedResponses2["BBBB"]), 1, t)

	lines3 := groupedResponses2["W"]
	groupedResponses3 := buildGroupedResponses(lines3, "CEBF")

	test.AssertDeepEqual(len(groupedResponses3["WWWW"]), 1, t)
	test.AssertDeepEqual(len(groupedResponses3["WWW"]), 4, t)
	test.AssertDeepEqual(len(groupedResponses3["WW"]), 6, t)
	test.AssertDeepEqual(len(groupedResponses3["W"]), 1, t)
	test.AssertDeepEqual(len(groupedResponses3["!"]), 0, t)

	test.AssertDeepEqual(len(groupedResponses3["BWWW"]), 3, t)
	test.AssertDeepEqual(len(groupedResponses3["BWW"]), 7, t)
	test.AssertDeepEqual(len(groupedResponses3["BW"]), 5, t)
	test.AssertDeepEqual(len(groupedResponses3["B"]), 1, t)

	test.AssertDeepEqual(len(groupedResponses3["BBWW"]), 3, t)
	test.AssertDeepEqual(len(groupedResponses3["BBW"]), 5, t)
	test.AssertDeepEqual(len(groupedResponses3["BB"]), 3, t)

	test.AssertDeepEqual(len(groupedResponses3["BBBW"]), 0, t)
	test.AssertDeepEqual(len(groupedResponses3["BBB"]), 4, t)

	test.AssertDeepEqual(len(groupedResponses3["BBBB"]), 1, t)

	lines4 := groupedResponses3["BWW"]
	groupedResponses4 := buildGroupedResponses(lines4, "ADFB")

	test.AssertDeepEqual(len(groupedResponses4["WWWW"]), 0, t)
	test.AssertDeepEqual(len(groupedResponses4["WWW"]), 1, t)
	test.AssertDeepEqual(len(groupedResponses4["WW"]), 1, t)
	test.AssertDeepEqual(len(groupedResponses4["W"]), 0, t)
	test.AssertDeepEqual(len(groupedResponses4["!"]), 0, t)

	test.AssertDeepEqual(len(groupedResponses4["BWWW"]), 0, t)
	test.AssertDeepEqual(len(groupedResponses4["BWW"]), 1, t)
	test.AssertDeepEqual(len(groupedResponses4["BW"]), 1, t)
	test.AssertDeepEqual(len(groupedResponses4["B"]), 1, t)

	test.AssertDeepEqual(len(groupedResponses4["BBWW"]), 0, t)
	test.AssertDeepEqual(len(groupedResponses4["BBW"]), 1, t)
	test.AssertDeepEqual(len(groupedResponses4["BB"]), 1, t)

	test.AssertDeepEqual(len(groupedResponses4["BBBW"]), 0, t)
	test.AssertDeepEqual(len(groupedResponses4["BBB"]), 0, t)

	test.AssertDeepEqual(len(groupedResponses4["BBBB"]), 0, t)

	test.AssertDeepEqual(groupedResponses4["BW"][0], "CFCB", t)
}

func buildGroupedResponses(lines []string, action string) map[string][]string {
	groupedResponses := make(map[string][]string)
	for _, line := range lines {
		world := game.NewWorld(line)
		reaction := Act(game.InitGamePlay(0.0), world, game.NewAction(action))
		responseString := reaction.Response.String()
		if _, ok := groupedResponses[responseString]; !ok {
			groupedResponses[responseString] = []string{}
		}
		groupedResponses[responseString] = append(groupedResponses[responseString], world.String())
	}

	return groupedResponses
}
