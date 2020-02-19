package game

// ###########################
// ### THE PLAYER ###
// ###########################

type Perspective struct {
	Knowledge       Knowledge
	PossibleWorlds  PossibleWorlds
	PossibleActions PossibleActions
	Agreement       Agreement
	Motivation      Motivation
}

func (p *Perspective) UpdateKnowledge(reaction Reaction) {
	p.Knowledge = append(p.Knowledge, reaction)
}

type Knowledge []Reaction

type PossibleWorlds []PossibleWorld
type PossibleWorld struct {
	World  World
	Weight float32
}

type PossibleActions []Action

type Motivation struct {
	Reward  float64
	Entropy float64
	Time    float64
}

func ExploreReactions(possibleWorlds PossibleWorlds, actions PossibleActions, agreement Agreement) ReactionCounter {
	reactionCounts := ReactionCounter{}
	for _, action := range actions {
		for _, possibleWorld := range possibleWorlds {
			reaction := agreement(possibleWorld.World, action)
			reactionCounts.Increment(reaction)
		}
	}
	// fmt.Printf("\nREACTIONCOUNTS: %+v", reactionCounts)
	return reactionCounts
}
