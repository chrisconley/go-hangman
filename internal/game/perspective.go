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
	// fmt.Printf("EXPLORING WITH WORLDS: %+v\n", possibleWorlds.ToNums())
	for _, action := range actions {
		hash := BuildActionWorldHash(action, possibleWorlds)
		if _, ok := BBPossibleWorldsStore[hash]; !ok {
			counts := utils.NewGroupedIncrementer()
			rewards := utils.NewGroupedIncrementer()
			responseWorlds := make(map[string]BPossibleWorlds)
			worldIndices := possibleWorlds.ToNums()
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
			// // if hash == uint64(16657282531674121003) {
			// // 	LogCacheWorlds(
			// // 		"CURENT WORLDS",
			// // 		hash, action, "n/a",
			// // 		possibleWorlds,
			// // 	)

			// // 	for response, pw := range responseWorlds {
			// // 		LogCacheWorlds(
			// // 			"SETTING CACHE",
			// // 			hash, action, response,
			// // 			pw,
			// // 		)
			// // 	}

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
}
