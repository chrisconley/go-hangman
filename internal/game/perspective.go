package game

// ###########################
// ### THE PLAYER ###
// ###########################

type Perspective struct {
	Motivation Motivation
	Agreement  Agreement
}

type Motivation struct {
	Reward    float64
	Entropy   float64
	MiniMax   float64
	MostParts float64
	Time      float64
}

func NewPerspective(
	motivation Motivation,
	agreement Agreement,
) Perspective {
	perspective := Perspective{}
	perspective.Motivation = motivation
	perspective.Agreement = agreement
	return perspective
}
