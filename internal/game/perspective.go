package game

import (
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/mitchellh/hashstructure"
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
	reactionExplorer := NewReactionExplorer()
	for _, action := range actions {
		counts := utils.NewGroupedIncrementer()
		rewards := utils.NewGroupedIncrementer()
		responseWorlds := make(map[string]BPossibleWorlds)
		for index, world := range BUniverse {
			if possibleWorlds.GetBit(uint64(index)) {
				reaction := agreement(World(world.Name), action)
				response := string(reaction.Response)
				if _, ok := responseWorlds[response]; !ok {
					responseWorlds[response] = NewPossibleWorlds(uint64(len(BUniverse)))
				}
				counts.Increment(response, 1.0)
				rewards.Increment(response, reaction.Reward)
				responseWorlds[response].Bits.SetBit(world.Index)
			}

			// reactionExplorer.Add(
			// 	reaction.Action,
			// 	reaction.Response,
			// 	reaction.Reward,
			// 	World(world.Name),
			// )
		}
		actionHash, err := hashstructure.Hash(action, nil)
		if err != nil {
			panic(err)
		}
		hash, err := hashstructure.Hash([]uint64{actionHash, possibleWorlds.GetHash()}, nil)
		if err != nil {
			panic(err)
		}
		BCountsStore[hash] = counts
		BRewardsStore[hash] = rewards
		BBPossibleWorldsStore[hash] = responseWorlds
		// for response, pw := range responseWorlds {
		// 	remainingWorldIndices := pw.Bits.ToNums()
		// 	fmt.Printf(
		// 		"remaining world indices: %d, %s, %s, %v\n",
		// 		hash,
		// 		action,
		// 		response,
		// 		remainingWorldIndices,
		// 	)

		// }
	}
	// fmt.Printf("\nREACTIONCOUNTS: %+v", reactionCounts)
	return reactionExplorer
}
