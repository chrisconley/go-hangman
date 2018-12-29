package game

import (
	"fmt"
)

// ### The system ###
type World string
type Action string
type Response string
type Reaction struct {
	Response Response
	Reward   float64
}
type Agreement func(w World, action Action) (reaction Reaction)

// ### The protagonist, player
// TODO: How can this be simplified? There are too many interfaces
type PossibleWorld struct {
	World  World
	Weight float32
}
type PossibleWorlds []PossibleWorld

type PossibleActions []Action

type Motivation struct {
	reward float32
	info   float32
	time   float32
}

type Perspective struct {
	possibleWorlds  PossibleWorlds
	possibleActions PossibleActions
	motivation      Motivation
}

type Predictions map[Action]PredictedActionUtility

type PredictedActionUtility struct {
	Action Action
	//reward float32
	Entropy float64
	//time   float32
}

func (p Predictions) GetEntropy(action string) float64 {
	return p[Action(action)].Entropy
}

// ### Process ###
func Explore(perspective Perspective, agreement Agreement) (predictions Predictions) {
	return Predictions{}
}

func Decide(predictions Predictions, motivation Motivation) (action Action) {
	return Action("a")
}

func Act(action Action, agreement Agreement, world World) (reaction Reaction) {
	return Reaction{"!", 0.5}
}

func UpdatePerspective(reaction Reaction, perspective Perspective) (newPerspective Perspective) {
	return perspective
}

// ### Process implementation ###
func PlaceholderHangmanAgreement(w World, d Action) (r Reaction) { return Reaction{"hi", 1.0} }
func main() {
	truth := World("abc")
	perspective := Perspective{
		possibleWorlds:  PossibleWorlds{PossibleWorld{"abc", 0.5}, PossibleWorld{"cab", 0.5}},
		possibleActions: PossibleActions{"a", "b", "c"},
		motivation:      Motivation{0.5, 0.25, 0.25},
	}

	predictions := Explore(perspective, PlaceholderHangmanAgreement)
	action := Decide(predictions, perspective.motivation)
	reaction := Act(action, PlaceholderHangmanAgreement, truth)
	newPerspective := UpdatePerspective(reaction, perspective)
	fmt.Printf("\n%+v", newPerspective)
}
