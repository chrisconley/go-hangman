package game

import (
	"github.com/chrisconley/go-hangman/internal/utils"
)

// ###########################
// ### THE PLAYER ###
// ###########################

type Perspective struct {
	Knowledge        Knowledge
	PossibleWorlds   PossibleWorlds
	AvailableActions AvailableActions
	Agreement        Agreement
	ActionFilter     ActionFilter
	Motivation       Motivation
	BPossibleWorlds  BPossibleWorlds
}

func NewPerspective(
	possibleWorlds PossibleWorlds,
	motivation Motivation,
	bpw BPossibleWorlds,
	agreement Agreement,
	actionFilter ActionFilter,
) *Perspective {
	perspective := &Perspective{}
	perspective.Motivation = motivation
	perspective.Knowledge = Knowledge{}
	perspective.PossibleWorlds = possibleWorlds
	perspective.BPossibleWorlds = bpw
	perspective.AvailableActions = BuildAvailableActions(possibleWorlds)
	perspective.Agreement = agreement
	perspective.ActionFilter = actionFilter
	return perspective
}

func (p *Perspective) UpdateKnowledge(reaction Reaction) {
	p.Knowledge = append(p.Knowledge, reaction)
}

func (p *Perspective) PlausibleActions(remainingWorlds PossibleWorlds) []Action {
	return p.ActionFilter(remainingWorlds, p.AvailableActions)
}

type Knowledge []Reaction

type Motivation struct {
	Reward  float64
	Entropy float64
	Time    float64
}

func ExploreReactions(possibleWorlds BPossibleWorlds, actions []Action, agreement Agreement) {
	for _, action := range actions {
		hash := BuildActionWorldHash(action, possibleWorlds)
		if _, ok := BBPossibleWorldsStore[hash]; !ok {
			counts, rewards, responseWorlds := ExploreReactionsForAction(
				possibleWorlds,
				action,
				agreement,
			)
			BCountsStore[hash] = counts
			BRewardsStore[hash] = rewards
			BBPossibleWorldsStore[hash] = responseWorlds
		}
	}
}

func ExploreReactionsForAction(
	possibleWorlds BPossibleWorlds,
	action Action,
	agreement Agreement,
) (
	utils.GroupedIncrementer,
	utils.GroupedIncrementer,
	map[string]BPossibleWorlds,
) {
	counts := utils.NewGroupedIncrementer()
	rewards := utils.NewGroupedIncrementer()
	responseWorlds := make(map[string]BPossibleWorlds)
	worldIndices := possibleWorlds.ToNums()
	for _, worldIndex := range worldIndices {
		world := BUniverse[worldIndex]
		reaction := agreement(World(world.Name), action)
		response := string(reaction.Response)
		if _, ok := responseWorlds[response]; !ok {
			responseWorlds[response] = NewPossibleWorlds(len(BUniverse))
		}
		counts.Increment(response, 1.0)
		rewards.Increment(response, reaction.Reward)
		worlds := responseWorlds[response]
		worlds.SetBit(world.Index)

	}
	return counts, rewards, responseWorlds
}
