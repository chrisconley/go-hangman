package game

import (
	"fmt"
	"math"
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
	Reward  float64
	Entropy float64
	Time    float64
}

type Perspective struct {
	possibleWorlds  PossibleWorlds
	possibleActions PossibleActions
	motivation      Motivation
}

type Predictions map[Action]PredictedActionUtility

type PredictedActionUtility struct {
	Action  Action
	Reward  float64
	Entropy float64
	//time   float32
}

func (p Predictions) GetEntropy(action string) float64 {
	return p[Action(action)].Entropy
}

func (p Predictions) AddPrediction(action Action, prediction PredictedActionUtility) {
	p[action] = prediction
}

// ### Process ###
func Explore(perspective Perspective, agreement Agreement) (predictions Predictions) {
	return Predictions{}
}

func Decide(predictions Predictions, motivation Motivation) (action Action) {
	maxWeight := 0.0
	maxAction := Action("")
	for _, prediction := range predictions {
		weight := math.Pow(prediction.Entropy, motivation.Entropy) * math.Pow(prediction.Reward, motivation.Reward)
		if weight > maxWeight {
			maxWeight = weight
			maxAction = prediction.Action
		}
	}
	return maxAction
}

func Act(action Action, agreement Agreement, truth World) (reaction Reaction) {
	return agreement(truth, action)
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
	// Where does the Goal come into play? Is that implicit in the Agreement?
	predictions := Explore(perspective, PlaceholderHangmanAgreement)
	action := Decide(predictions, perspective.motivation)
	reaction := Act(action, PlaceholderHangmanAgreement, truth)
	newPerspective := UpdatePerspective(reaction, perspective)
	// TODO: Change this to Observe, which captures reaction and keeps world state
	fmt.Printf("\n%+v", newPerspective)
}
