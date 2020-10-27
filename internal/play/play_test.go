package play

import (
	"fmt"
	"sort"
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/utils"

	"github.com/chrisconley/go-hangman/test"
)

func TestAct(t *testing.T) {
	truth := game.NewWorld("abc")
	actual := Act(game.InitGamePlay(0.0), game.NewAction("a"), game.BuildTestAct(0.222), truth)
	test.AssertDeepEqual(
		actual,
		game.Reaction{
			Action:   game.NewAction("a"),
			Response: game.NewResponse("a"),
			Reward:   0.222,
		},
		t,
	)
}

func TestObserve(t *testing.T) {
	t.Skip("TODO: Fix test setup")
}

type ResponseCounts struct {
	counts            map[string]utils.GroupedIncrementer
	PossibleResponses map[string]struct{}
}

func NewResponseCounts(
	gameplay game.GamePlay,
	remainingWorlds []game.World,
	actions []game.Action,
	act game.Act,
) ResponseCounts {
	rc := ResponseCounts{}
	rc.counts = make(map[string]utils.GroupedIncrementer, len(actions))
	rc.PossibleResponses = map[string]struct{}{}
	for _, action := range actions {
		for _, world := range remainingWorlds {
			reaction := act(gameplay, world, action)
			actionString := action.String()
			responseString := reaction.Response.String()
			if _, ok := rc.counts[actionString]; !ok {
				rc.counts[actionString] = utils.NewGroupedIncrementer()
			}
			c := rc.counts[actionString]
			c.Increment(responseString, 1.0)
			rc.counts[actionString] = c
			rc.PossibleResponses[responseString] = struct{}{}
		}
	}
	return rc
}

func (rc ResponseCounts) Values() map[string]utils.GroupedIncrementer {
	return rc.counts
}

func (rc ResponseCounts) Get(actionString string) utils.GroupedIncrementer {
	return rc.counts[actionString]
}

func (rc *ResponseCounts) Set(actionString string, gi utils.GroupedIncrementer) {
	rc.counts[actionString] = gi
}

type GroupedCounter struct {
	counts  utils.GroupedIncrementer
	Actions []string
}

func (rc ResponseCounts) GroupByCountHash(groups []string) map[uint64]GroupedCounter {
	groupedCounter := map[uint64]GroupedCounter{}
	for actionString, groupedIncrementer := range rc.counts {
		hash := groupedIncrementer.GetHash(groups)
		counter := groupedCounter[hash]
		counter.counts = groupedIncrementer
		counter.Actions = append(counter.Actions, actionString)
		groupedCounter[hash] = counter
	}
	return groupedCounter
}

type PotentialSymmetries struct {
	GroupedActions    [][]string
	ResponseCounts    ResponseCounts
	possibleResponses []string
	worldActions      []game.Action
}

func (ps *PotentialSymmetries) Print() {
	for _, worldAction := range ps.worldActions {
		counts := ps.ResponseCounts.Get(worldAction.String())
		fmt.Printf("%s: ", worldAction)
		for _, response := range ps.possibleResponses {
			fmt.Printf("%s: %.0f ", response, counts.Get(response))
		}
		ps.ResponseCounts.Set(worldAction.String(), counts)
		fmt.Printf("\n")
	}

}

func NewPotentialSymmetries(reactions []game.Reaction) PotentialSymmetries {
	ps := PotentialSymmetries{}
	ps.worldActions = game.BUniverse.Actions.WorldActions()
	possibleWorlds := game.BUniverse.PossibleWorlds
	gameplay := game.InitGamePlay(possibleWorlds.GetHash())
	for _, reaction := range reactions {
		_, possibleWorlds = Observe(gameplay, reaction, mastermind.Act, possibleWorlds)

	}
	worlds := game.PossibleWorldsToWorlds(possibleWorlds)
	fmt.Printf("Len: %d | %+v\n", possibleWorlds.Len(), worlds)
	ps.ResponseCounts = NewResponseCounts(
		gameplay,
		worlds,
		ps.worldActions,
		mastermind.Act,
	)

	for response := range ps.ResponseCounts.PossibleResponses {
		ps.possibleResponses = append(ps.possibleResponses, response)
	}
	fmt.Println(ps.ResponseCounts.PossibleResponses)
	for _, worldAction := range ps.worldActions {
		counts := ps.ResponseCounts.Get(worldAction.String())
		for _, response := range ps.possibleResponses {
			counts.Increment(response, 0)
		}
		ps.ResponseCounts.Set(worldAction.String(), counts)
	}
	return ps
}

func findPotentialSymmetries(
	reactions []game.Reaction,
) PotentialSymmetries {
	ps := NewPotentialSymmetries(reactions)

	groupedActions := [][]string{}
	for _, groupedCounter := range ps.ResponseCounts.GroupByCountHash(ps.possibleResponses) {
		sort.Strings(groupedCounter.Actions)
		groupedActions = append(groupedActions, groupedCounter.Actions)

	}
	sort.Slice(groupedActions, func(i, j int) bool {
		return groupedActions[i][0] < groupedActions[j][0]
	})
	ps.GroupedActions = groupedActions
	return ps
}
func TestPlayMastermindSymmetry2(t *testing.T) {
	t.Skip("hi")
	lines := []string{
		"AA", "AB", "AC", "AD",
		"BA", "BB", "BC", "BD",
		"CA", "CB", "CC", "CD",
		"DA", "DB", "DC", "DD",
	}
	game.BUniverse = game.InitUniverse(lines)

	reactions := []game.Reaction{
		game.NewReaction(game.NewAction("AB"), "B", 0.0),
	}
	potentialSymmetries := findPotentialSymmetries(
		reactions,
	)

	potentialSymmetries.Print()
	fmt.Printf("%+v\n", potentialSymmetries.GroupedActions)
	test.AssertDeepEqual(2, 1, t)
}

func TestPlayMastermindSymmetry3(t *testing.T) {
	t.Skip("hi")
	lines := utils.ReadLinesFromFilePath("../../fixtures/mastermind-ABCD-3.txt")
	// lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")
	game.BUniverse = game.InitUniverse(lines)
	reactions := []game.Reaction{
		game.NewReaction(game.NewAction("ABC"), "WW", 0.0),
		game.NewReaction(game.NewAction("ABD"), "WW", 0.0),
		// game.NewReaction(game.NewAction("ABB"), "BWW", 0.0),
	}

	// reactions := []game.Reaction{
	// 	game.NewReaction(game.NewAction("AABB"), "B", 0.0),
	// 	// game.NewReaction(game.NewAction("ACCD"), "B", 0.0),
	// }

	// [AAB ACA     BBA     BCC     CAC CBB    ]
	// [AAB     ACC BBA BCB     CAA         CBC]
	potentialSymmetries := findPotentialSymmetries(
		reactions,
	)

	potentialSymmetries.Print()
	fmt.Printf("Len: %d | %+v\n", len(potentialSymmetries.GroupedActions), potentialSymmetries.GroupedActions)
	test.AssertDeepEqual(2, 1, t)
}

func TestPlayMastermindSymmetry(t *testing.T) {
	t.Skip("hi")
	lines := []string{
		"AA", "AB", "AC", "AD",
		"BA", "BB", "BC", "BD",
		"CA", "CB", "CC", "CD",
		"DA", "DB", "DC", "DD",
	}
	game.BUniverse = game.InitUniverse(lines)
	agreement := mastermind.Agreement
	// perspective := game.NewPerspective(
	// 	game.Motivation{Reward: 0.5, Entropy: 0.25, Time: 0.25},
	// 	agreement,
	// )
	possibleWorlds := game.BUniverse.PossibleWorlds
	gameplay := game.InitGamePlay(possibleWorlds.GetHash())
	predictions := Explore(gameplay, agreement, possibleWorlds)
	possibleResponses := []string{"BB", "B", "WW", "W", "!"}
	newPredictions := getNewPredictions(
		predictions,
		possibleResponses,
		gameplay,
		mastermind.Act,
		mastermind.Agreement,
		possibleWorlds,
	)
	for actionString, predictionsSlice := range newPredictions {
		for _, predictions := range predictionsSlice {
			fmt.Printf("%s %+v\n", actionString, predictions.PredictionsMap.Values)
		}
	}
	test.AssertDeepEqual(newPredictions, []game.Prediction{}, t)
}

func getNewPredictions(
	predictions game.Predictions,
	possibleResponses []string,
	gameplay game.GamePlay,
	act game.Act,
	agreement game.Agreement,
	possibleWorlds game.PossibleWorlds,
) map[string][]game.Predictions {
	newPredictions := map[string][]game.Predictions{}
	values := predictions.PredictionsMap.Values
	for _, prediction := range values {
		action := prediction.Action
		for _, responseString := range possibleResponses {
			_, npw := Observe(
				gameplay,
				game.NewReaction(action, responseString, 0.0),
				act,
				possibleWorlds,
			)
			fmt.Println(npw.Len())
			fmt.Println(Explore(gameplay, agreement, npw))
			newPredictions[action.String()] = append(
				newPredictions[action.String()],
				Explore(gameplay, agreement, npw),
			)

		}

	}
	return newPredictions

}

func TestPlayBHangman(t *testing.T) {
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
	game.BUniverse = game.InitUniverse(lines)

	for _, line := range lines {
		perspective := game.NewPerspective(
			game.Motivation{Reward: 0.5, Entropy: 0.25, Time: 0.25},
			hangman.Agreement,
		)
		truth := game.NewWorld(line)
		_, actual := PlayGame(truth, perspective)
		test.AssertDeepEqual(actual, truth, t)
	}
}

func TestPlayMastermind(t *testing.T) {
	lines := []string{
		"111",
		"112",
		"121",
		"211",
	}
	game.BUniverse = game.InitUniverse(lines)

	for _, line := range lines {
		perspective := game.NewPerspective(
			game.Motivation{Reward: 0.5, Entropy: 0.25, Time: 0.25},
			mastermind.Agreement,
		)

		truth := game.NewWorld(line)
		_, actual := PlayGame(truth, perspective)
		test.AssertDeepEqual(actual, truth, t)
	}
}
func TestIntegrationPlayKnuthMastermind(t *testing.T) {
	t.Skip("Skipping TestIntegrationPlayKnuthMastermind since it's no longer correct with mastermind action pruing.")
	// In Knuth's paper, the second guess is ACDD since it leaves at most 44 remaining choices
	// and is lexigraphically first out of other choices that also minimize the max remaining choices.
	// Since the first guess is AABB, the remaining "free colors" are C, D, E, and F.
	// When pruning symmetrical actions, ACDD and ACCD are symetrical and ACCD is used because
	// it is lexigraphically before ACDD. ACCD leaves at most 51 remaining choices and is therefore
	// not selected. Instead CCBD is selected with 44 remaining choices.
	// Pruning symmetrical actions results in 5800 total turns instead of 5801.

	lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")
	game.BUniverse = game.InitUniverse(lines)

	perspective := game.NewPerspective(
		game.Motivation{Entropy: 0.00, MiniMax: 1.0},
		mastermind.Agreement,
	)

	truth := game.NewWorld("CFCB")
	gameplay, _ := PlayGame(truth, perspective)

	test.AssertDeepEqual(gameplay.Reactions[0].Action, game.NewAction("AABB"), t)
	test.AssertDeepEqual(gameplay.Reactions[0].Response, game.NewResponse("B"), t)

	test.AssertDeepEqual(gameplay.Reactions[1].Action, game.NewAction("ACDD"), t)
	test.AssertDeepEqual(gameplay.Reactions[1].Response, game.NewResponse("W"), t)

	test.AssertDeepEqual(gameplay.Reactions[2].Action, game.NewAction("CEBF"), t)
	test.AssertDeepEqual(gameplay.Reactions[2].Response, game.NewResponse("BWW"), t)

	test.AssertDeepEqual(gameplay.Reactions[3].Action, game.NewAction("ADFB"), t)
	test.AssertDeepEqual(gameplay.Reactions[3].Response, game.NewResponse("BW"), t)
}

func TestIntegrationVerifyMiniMaxMastermind(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestIntegrationVerifyMiniMaxMastermind in short mode.")
	}
	lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")

	motivation := game.Motivation{MiniMax: 1.0}
	limit := len(lines)
	games := PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := CalculateStats(games)
	test.AssertDeepEqual(stats.GameCount, 1296, t)
	test.AssertDeepEqual(stats.MaxTurns, 5, t)
	// Note that Ville (Optimal Mastermind Strategy) shows that total turns is 5801.
	// Because of symmetrical action pruning, this implmentation reduces that by 1 to 5800.
	test.AssertDeepEqual(stats.TotalTurns, 5800, t)
	test.AssertAlmostEqual(stats.AverageTurns, 4.4753, 4, t)
}

func TestIntegrationVerifyMostPartsMastermind(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestIntegrationVerifyMiniMaxMastermind in short mode.")
	}
	lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")

	motivation := game.Motivation{MostParts: 1.0}
	limit := len(lines)
	games := PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := CalculateStats(games)
	test.AssertDeepEqual(stats.GameCount, 1296, t)
	test.AssertDeepEqual(stats.MaxTurns, 6, t)
	test.AssertDeepEqual(stats.TotalTurns, 5668, t)
	test.AssertAlmostEqual(stats.AverageTurns, 4.3735, 4, t)
}

func TestIntegrationVerifyEntropyMastermind(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestIntegrationVerifyEntropyMastermind in short mode.")
	}
	lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")

	motivation := game.Motivation{Entropy: 1.0}
	limit := len(lines)
	games := PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := CalculateStats(games)
	test.AssertDeepEqual(stats.GameCount, 1296, t)
	test.AssertDeepEqual(stats.MaxTurns, 6, t)
	test.AssertDeepEqual(stats.TotalTurns, 5723, t)
	test.AssertAlmostEqual(stats.AverageTurns, 4.4159, 4, t)
}

func TestIntegrationVerifyEntropyHangman(t *testing.T) {
	t.Skip("Skipping TestIntegrationVerifyEntropyHangman because it's hanging.")
	if testing.Short() {
		t.Skip("Skipping TestIntegrationVerifyEntropyHangman in short mode.")
	}
	lines := utils.ReadLinesFromFilePath("../../fixtures/nine-letter-words.txt")

	motivation := game.Motivation{Entropy: 1.0}
	limit := len(lines)
	games := PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := CalculateStats(games)
	test.AssertDeepEqual(stats.GameCount, 25011, t)
	test.AssertDeepEqual(stats.MaxTurns, 13, t)
	test.AssertDeepEqual(stats.TotalTurns, 146356, t)
	test.AssertAlmostEqual(stats.AverageTurns, 5.8529, 4, t)
}

func TestIntegrationVerifyMiniMaxHangman(t *testing.T) {
	t.Skip("Skipping TestIntegrationVerifyMiniMaxHangman because it's hanging.")
	if testing.Short() {
		t.Skip("Skipping TestIntegrationVerifyMiniMaxHangman in short mode.")
	}
	lines := utils.ReadLinesFromFilePath("../../fixtures/nine-letter-words.txt")

	motivation := game.Motivation{MiniMax: 1.0}
	limit := len(lines)
	games := PlayGames(mastermind.Agreement, motivation, lines, limit)
	stats := CalculateStats(games)
	test.AssertDeepEqual(stats.GameCount, 25011, t)
	test.AssertDeepEqual(stats.MaxTurns, 12, t)
	test.AssertDeepEqual(stats.TotalTurns, 149552, t)
	test.AssertAlmostEqual(stats.AverageTurns, 5.9794, 4, t)
}

func TestPlayGames(t *testing.T) {
	lines := utils.ReadLinesFromFilePath("../../fixtures/nine-letter-words.txt")
	lines = lines[:100]

	motivation := game.Motivation{Reward: 0.0, Entropy: 1.0, Time: 0.0}
	games := PlayGames(mastermind.Agreement, motivation, lines, 3)
	test.AssertDeepEqual(games[0].TurnCount, 2, t)
	test.AssertDeepEqual(games[1].TurnCount, 3, t)
	test.AssertDeepEqual(games[2].TurnCount, 3, t)
}

func TestBuildTrialsTwoDimensions(t *testing.T) {
	trials := BuildTrials(0.0, 1.0, 0.25, 2)
	test.AssertDeepEqual(len(trials), 5, t)
	test.AssertDeepEqual(trials[0], []float64{0.0, 1.0}, t)
	test.AssertDeepEqual(trials[1], []float64{0.25, 0.75}, t)
	test.AssertDeepEqual(trials[2], []float64{0.5, 0.5}, t)
	test.AssertDeepEqual(trials[3], []float64{0.75, 0.25}, t)
	test.AssertDeepEqual(trials[4], []float64{1.0, 0.0}, t)
}
func TestBuildTrialsThreeDimensions(t *testing.T) {
	trials := BuildTrials(0.0, 1.0, 0.25, 3)
	test.AssertDeepEqual(len(trials), 15, t)
	test.AssertDeepEqual(trials[0], []float64{0.0, 0.0, 1.0}, t)
	test.AssertDeepEqual(trials[1], []float64{0.0, 0.25, 0.75}, t)
	test.AssertDeepEqual(trials[2], []float64{0.0, 0.5, 0.5}, t)
	test.AssertDeepEqual(trials[3], []float64{0.0, 0.75, 0.25}, t)
	test.AssertDeepEqual(trials[4], []float64{0.0, 1.0, 0.0}, t)

	test.AssertDeepEqual(trials[5], []float64{0.25, 0.0, 0.75}, t)
	test.AssertDeepEqual(trials[6], []float64{0.25, 0.25, 0.5}, t)
	test.AssertDeepEqual(trials[7], []float64{0.25, 0.5, 0.25}, t)
	test.AssertDeepEqual(trials[8], []float64{0.25, 0.75, 0.0}, t)

	test.AssertDeepEqual(trials[9], []float64{0.5, 0.0, 0.5}, t)
	test.AssertDeepEqual(trials[10], []float64{0.5, 0.25, 0.25}, t)
	test.AssertDeepEqual(trials[11], []float64{0.5, 0.5, 0.0}, t)

	test.AssertDeepEqual(trials[12], []float64{0.75, 0.0, 0.25}, t)
	test.AssertDeepEqual(trials[13], []float64{0.75, 0.25, 0.0}, t)

	test.AssertDeepEqual(trials[14], []float64{1.0, 0.0, 0.0}, t)
}
func TestBuildTrialsFourDimensions(t *testing.T) {
	trials := BuildTrials(0.0, 1.0, 0.25, 4)
	test.AssertDeepEqual(len(trials), 35, t)
	test.AssertDeepEqual(trials[0], []float64{0.0, 0.0, 0.0, 1.0}, t)
	test.AssertDeepEqual(trials[14], []float64{0.0, 1.0, 0.0, 0.0}, t)
	test.AssertDeepEqual(trials[15], []float64{0.25, 0.0, 0.0, 0.75}, t)
	test.AssertDeepEqual(trials[24], []float64{0.25, 0.75, 0.0, 0.0}, t)
	test.AssertDeepEqual(trials[25], []float64{0.5, 0.0, 0.0, 0.5}, t)
	test.AssertDeepEqual(trials[30], []float64{0.5, 0.5, 0.0, 0.0}, t)
	test.AssertDeepEqual(trials[31], []float64{0.75, 0.0, 0.0, 0.25}, t)
	test.AssertDeepEqual(trials[33], []float64{0.75, 0.25, 0.0, 0.0}, t)
	test.AssertDeepEqual(trials[34], []float64{1.0, 0.0, 0.0, 0.0}, t)
}

func TestBuildTrialsFourDimensionSmallerStep(t *testing.T) {
	trials := BuildTrials(0.0, 1.0, 0.2, 4)
	test.AssertDeepEqual(len(trials), 56, t)

	trials = BuildTrials(0.0, 1.0, 0.1, 3)
	test.AssertDeepEqual(len(trials), 66, t)

	trials = BuildTrials(0.0, 1.0, 0.05, 3)
	test.AssertDeepEqual(len(trials), 231, t)

	trials = BuildTrials(0.0, 1.0, 0.01, 3)
	test.AssertDeepEqual(len(trials), 5151, t)
}
