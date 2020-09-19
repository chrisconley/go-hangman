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
	ReactionExplorer reactionExplorer
	BPossibleWorlds  BPossibleWorlds
}

func NewPerspective(possibleWorlds PossibleWorlds, motivation Motivation) *Perspective {
	perspective := &Perspective{}
	perspective.Motivation = motivation
	perspective.Knowledge = Knowledge{}
	perspective.PossibleWorlds = possibleWorlds
	perspective.AvailableActions = BuildAvailableActions(possibleWorlds)
	return perspective
}

func (p *Perspective) UpdateKnowledge(reaction Reaction) {
	p.Knowledge = append(p.Knowledge, reaction)
}

func (p *Perspective) UpdateReactionExplorer(reactionExplorer reactionExplorer) {
	p.ReactionExplorer = reactionExplorer
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

func ExploreReactions(possibleWorlds BPossibleWorlds, actions []Action, agreement Agreement) reactionExplorer {
	// fmt.Printf("EXPLORING WITH WORLDS: %+v\n", possibleWorlds.Bits.ToNums())
	for _, action := range actions {
		hash := BuildActionWorldHash(action, possibleWorlds)
		if _, ok := BBPossibleWorldsStore[hash]; !ok {
			counts := utils.NewGroupedIncrementer()
			rewards := utils.NewGroupedIncrementer()
			responseWorlds := make(map[string]BPossibleWorlds)
			worldIndices := possibleWorlds.Bits.ToNums()
			for _, worldIndex := range worldIndices {
				world := BUniverse[worldIndex]
				reaction := agreement(World(world.Name), action)
				response := string(reaction.Response)
				if _, ok := responseWorlds[response]; !ok {
					responseWorlds[response] = NewPossibleWorlds(uint64(len(BUniverse)))
				}
				counts.Increment(response, 1.0)
				rewards.Increment(response, reaction.Reward)
				responseWorlds[response].Bits.SetBit(world.Index)

			}

			BCountsStore[hash] = counts
			BRewardsStore[hash] = rewards
			BBPossibleWorldsStore[hash] = responseWorlds
			// for response, pw := range responseWorlds {
			// 	LogCacheWorlds(
			// 		"SETTING CACHE",
			// 		hash, action, response,
			// 		pw,
			// 	)
			// }
		} else {
			// do nothing because we already have what we need
			// responseWorlds := BBPossibleWorldsStore[hash]
			// for response, pw := range responseWorlds {
			// 	LogCacheWorlds(
			// 		"USING CACHE",
			// 		hash, action, response,
			// 		pw,
			// 	)
			// }
		}
	}
	// fmt.Printf("\nREACTIONCOUNTS: %+v", reactionCounts)
	reactionExplorer := NewReactionExplorer()
	return reactionExplorer
}
