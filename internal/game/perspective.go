package game

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

func ExploreReactions(possibleWorlds PossibleWorlds, actions []Action, agreement Agreement) reactionExplorer {
	reactionExplorer := NewReactionExplorer()
	for _, action := range actions {
		for _, possibleWorld := range possibleWorlds.Worlds {
			reaction := agreement(possibleWorld, action)
			reactionExplorer.Add(
				reaction.Action,
				reaction.Response,
				reaction.Reward,
				possibleWorld,
			)
		}
	}
	// fmt.Printf("\nREACTIONCOUNTS: %+v", reactionCounts)
	return reactionExplorer
}
