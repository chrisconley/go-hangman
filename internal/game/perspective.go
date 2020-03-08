package game

// ###########################
// ### THE PLAYER ###
// ###########################

type Perspective struct {
	Knowledge        Knowledge
	PossibleWorlds   Worlds
	AvailableActions AvailableActions
	Agreement        Agreement
	Motivation       Motivation
	ReactionExplorer reactionExplorer
}

// func NewPerspective(worlds Worlds, motivation Motivation, game string) *Perspective {
// 	perspective := &Perspective{}
// 	perspective.Knowledge = Knowledge{}
// 	perspective.PossibleWorlds = worlds
// 	perspective.AvailableActions = BuildAvailableActions(worlds)
// 	if game == "hangman" {
// 		perspective.Agreement = hangman.HangmanAgreement
// 	} else if game == "mastermind" {
// 		perspective.Agreement = mastermind.Agreement
// 	}
// }

func (p *Perspective) UpdateKnowledge(reaction Reaction) {
	p.Knowledge = append(p.Knowledge, reaction)
}

func (p *Perspective) UpdateReactionExplorer(reactionExplorer reactionExplorer) {
	p.ReactionExplorer = reactionExplorer
}

type Knowledge []Reaction

type Worlds []World

type Motivation struct {
	Reward  float64
	Entropy float64
	Time    float64
}

func ExploreReactions(possibleWorlds Worlds, actions AvailableActions, agreement Agreement) reactionExplorer {
	reactionExplorer := NewReactionExplorer()
	for _, action := range actions.All() {
		for _, possibleWorld := range possibleWorlds {
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
