package game

// ###########################
// ### THE PLAYER ###
// ###########################

type Perspective struct {
	Knowledge        Knowledge
	PossibleWorlds   PossibleWorlds
	PossibleActions  PossibleActions
	Agreement        Agreement
	Motivation       Motivation
	ReactionExplorer reactionExplorer
}

func (p *Perspective) UpdateKnowledge(reaction Reaction) {
	p.Knowledge = append(p.Knowledge, reaction)
}

func (p *Perspective) UpdateReactionExplorer(reactionExplorer reactionExplorer) {
	p.ReactionExplorer = reactionExplorer
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

func ExploreReactions(possibleWorlds PossibleWorlds, actions PossibleActions, agreement Agreement) reactionExplorer {
	reactionExplorer := NewReactionExplorer()
	for _, action := range actions {
		for _, possibleWorld := range possibleWorlds {
			reaction := agreement(possibleWorld.World, action)
			reactionExplorer.Add(
				reaction.Action,
				reaction.Response,
				reaction.Reward,
				possibleWorld.World,
			)
		}
	}
	// fmt.Printf("\nREACTIONCOUNTS: %+v", reactionCounts)
	return reactionExplorer
}
