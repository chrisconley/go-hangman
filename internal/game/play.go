package game

import "fmt"

// ### The true state of the world ###
// Special moniker for the solution, codeword, nature, reality
type Truth World

// ### The system ###
type World string
type Action string
type Reaction string
type Agreement func(w World, action Action) (reaction Reaction)

// ### The protagonist, player
// TODO: How can this be simplified? There are too many interfaces
type PossibleWorld struct {
	world  World
	weight float32
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
	action Action
	reward float32
	info   float32
	time   float32
}

// ### Process ###
func Explore(perspective Perspective, agreement Agreement) (predictions Predictions) {
	return Predictions{}
}

func Decide(predictions Predictions, motivation Motivation) (action Action) {
	return Action("a")
}

func Act(action Action, agreement Agreement, truth Truth) (reaction Reaction) {
	return Reaction("!")
}

func UpdatePerspective(reaction Reaction, perspective Perspective) (newPerspective Perspective) {
	return perspective
}

// ### Process implementation ###
func HangmanAgreement(w World, d Action) (r Reaction) { return Reaction("hi") }
func main() {
	truth := Truth("abc")
	perspective := Perspective{
		possibleWorlds:  PossibleWorlds{PossibleWorld{"abc", 0.5}, PossibleWorld{"cab", 0.5}},
		possibleActions: PossibleActions{"a", "b", "c"},
		motivation:      Motivation{0.5, 0.25, 0.25},
	}

	predictions := Explore(perspective, HangmanAgreement)
	action := Decide(predictions, perspective.motivation)
	reaction := Act(action, HangmanAgreement, truth)
	newPerspective := UpdatePerspective(reaction, perspective)
	fmt.Printf("\n%+v", newPerspective)
}
