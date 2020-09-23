package game

// ###########################
// ### THE PLAYER ###
// ###########################

type Perspective struct {
	Motivation Motivation
	BAgreement BAgreement
}

type Motivation struct {
	Reward  float64
	Entropy float64
	Time    float64
}

func NewPerspective(
	motivation Motivation,
	bagreement BAgreement,
) Perspective {
	perspective := Perspective{}
	perspective.Motivation = motivation
	perspective.BAgreement = bagreement
	return perspective
}
