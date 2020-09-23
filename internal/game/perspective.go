package game

// ###########################
// ### THE PLAYER ###
// ###########################

type Perspective struct {
	PossibleWorlds  PossibleWorlds
	Agreement       Agreement
	ActionFilter    ActionFilter
	Motivation      Motivation
	BPossibleWorlds BPossibleWorlds
}

func NewPerspective(
	motivation Motivation,
	agreement Agreement,
	actionFilter ActionFilter,
) Perspective {
	perspective := Perspective{}
	perspective.Motivation = motivation
	perspective.Agreement = agreement
	perspective.ActionFilter = actionFilter
	return perspective
}

type Motivation struct {
	Reward  float64
	Entropy float64
	Time    float64
}
