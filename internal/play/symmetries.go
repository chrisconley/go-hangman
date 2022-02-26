package play

import (
	"fmt"
	"sort"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/utils"
)

func CountUniqueCharacters(s string) int {
	set := make(map[rune]int)
	uniqueCount := 0
	for _, c := range s {
		if _, ok := set[c]; !ok {
			uniqueCount++
		}
		set[c]++
	}
	return uniqueCount
}

func CountHeldCharacters(g string, s string) int {
	heldCount := 0
	for index, char := range s {
		if g[index] == byte(char) {
			heldCount++
		}
	}
	return heldCount
}

func NewPotentialSymmetries(reactions []game.Reaction) ResponseCounts {
	worldActions := game.BUniverse.Actions.WorldActions()
	possibleWorlds := game.BUniverse.PossibleWorlds
	gameplay := game.InitGamePlay(possibleWorlds.GetHash())
	for _, reaction := range reactions {
		_, possibleWorlds = Observe(gameplay, reaction, mastermind.Act, possibleWorlds)
	}

	responseCounts := NewResponseCounts(
		gameplay,
		game.PossibleWorldsToWorlds(possibleWorlds),
		worldActions,
		mastermind.Act,
	)
	fmt.Printf("%+v\n", game.PossibleWorldsToWorlds(possibleWorlds))

	return responseCounts
}

type ResponseCounts struct {
	counts               map[string]utils.GroupedIncrementer
	possibleResponsesSet map[string]struct{}
}

func NewResponseCounts(
	gameplay game.GamePlay,
	remainingWorlds []game.World,
	actions []game.Action,
	act game.Act,
) ResponseCounts {
	rc := ResponseCounts{}
	rc.counts = make(map[string]utils.GroupedIncrementer, len(actions))
	rc.possibleResponsesSet = map[string]struct{}{}
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
			rc.possibleResponsesSet[responseString] = struct{}{}
		}
	}
	for _, worldAction := range actions {
		counts := rc.Get(worldAction.String())
		for _, response := range rc.SortedPossibleResponses() {
			counts.Increment(response, 0)
		}
		rc.Set(worldAction.String(), counts)
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

func (rc ResponseCounts) SortedPossibleResponses() []string {
	possibleResponses := []string{}
	for response := range rc.possibleResponsesSet {
		possibleResponses = append(possibleResponses, response)
	}

	sort.Slice(possibleResponses, func(i, j int) bool {
		return possibleResponses[i] < possibleResponses[j]
	})
	return possibleResponses
}

type GroupedCounter struct {
	counts  utils.GroupedIncrementer
	Actions []string
}

func (rc ResponseCounts) GroupByCountHash() map[uint64]GroupedCounter {
	groupedCounter := map[uint64]GroupedCounter{}
	for actionString, groupedIncrementer := range rc.counts {
		hash := groupedIncrementer.GetHash(rc.SortedPossibleResponses())
		counter := groupedCounter[hash]
		counter.counts = groupedIncrementer
		counter.Actions = append(counter.Actions, actionString)
		groupedCounter[hash] = counter
	}
	return groupedCounter
}

func (rc *ResponseCounts) GroupEquivalentActions() [][]string {
	groupedActions := [][]string{}
	for _, groupedCounter := range rc.GroupByCountHash() {
		sort.Strings(groupedCounter.Actions)
		groupedActions = append(groupedActions, groupedCounter.Actions)

	}
	sort.Slice(groupedActions, func(i, j int) bool {
		return groupedActions[i][0] < groupedActions[j][0]
	})
	return groupedActions
}

func (rc *ResponseCounts) ToString(actionString string, possibleResponses []string) string {
	result := ""
	counts := rc.Get(actionString)
	result += fmt.Sprintf("%s: ", actionString)
	for _, response := range possibleResponses {
		result += fmt.Sprintf("%s: %.0f\t", response, counts.Get(response))
	}
	return result
}
