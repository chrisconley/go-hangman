package game

// ###########################
// ### THE SYSTEM ###
// ###########################
type World string
type Action string
type Response string
type Agreement func(w World, action Action) (reaction Reaction)

// ###########################
// ### THE PLAYER ###
// ###########################

type Reaction struct {
	Action   Action
	Response Response
	Reward   float64
}

type Knowledge []Reaction

type PossibleWorld struct {
	World  World
	Weight float32
}
type PossibleWorlds []PossibleWorld

type PossibleActions []Action

type Motivation struct {
	Reward  float64
	Entropy float64
	Time    float64
}

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
