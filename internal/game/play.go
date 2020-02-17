package game

import (
	"fmt"
	// "math"
	// "regexp"
)

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
	Action Action
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
	Agreement		Agreement
	Motivation      Motivation
}

func (p *Perspective) UpdateKnowledge(reaction Reaction) {
	p.Knowledge = append(p.Knowledge, reaction)
}

// ###########################
// ### THE DECISION PROCESS ###
// ###########################
func Explore(perspective Perspective, agreement Agreement) (predictions Predictions) {
	reactions := ExploreReactions(perspective.PossibleWorlds, perspective.PossibleActions, agreement)
	predictions = BuildPredictions(reactions)
	return predictions
}

func Decide(predictions Predictions, motivation Motivation) (action Action) {
	maxWeight := 0.0
	maxAction := Action("")
	for _, prediction := range predictions {
		// fmt.Printf("\nPREDICTION: %+v", prediction)
		// weight := math.Pow(prediction.Entropy, motivation.Entropy) * math.Pow(prediction.Reward, motivation.Reward)
		weight := prediction.Entropy * motivation.Entropy + prediction.Reward * motivation.Reward
		// fmt.Printf("\nWEIGHT: %+v", weight)
		if weight > maxWeight {
			maxWeight = weight
			maxAction = prediction.Action
		}
	}
	// fmt.Printf("\nMAXACTION: %+v", maxAction)
	return maxAction
}

func Act(action Action, agreement Agreement, truth World) (reaction Reaction) {
	return agreement(truth, action)
}

func Observe(action Action, reaction Reaction, perspective *Perspective) {
	perspective.UpdateKnowledge(reaction)
	perspective.PossibleActions = UpdatePossibleActions(reaction, perspective.PossibleActions)

	perspective.PossibleWorlds = UpdatePossibleWorlds(reaction, perspective.PossibleWorlds)
	// fmt.Printf("\nPOSSIBLE ACTIONS: %+v\n", perspective.PossibleActions)
	// fmt.Printf("\nPOSSIBLE WORLDS: %+v\n", perspective.PossibleWorlds)

}

// ###########################
// ### THE GAME ###
// ###########################

func PlayGame(truth World, perspective *Perspective) (World) {
	for true {
		playTurn(truth, perspective)
		if len(perspective.PossibleWorlds) == 1 {
			break
		}
	}
	// fmt.Printf("\nGAME OVER: %+v\n", perspective.PossibleWorlds[0])
	return perspective.PossibleWorlds[0].World
}

func playTurn(truth World, perspective *Perspective) (reaction Reaction) {
	predictions := Explore(*perspective, HangmanAgreement)
	action := Decide(predictions, perspective.Motivation)
	// fmt.Printf("\nACTION: %+v", action)
	reaction = Act(action, HangmanAgreement, truth)
	// fmt.Printf("\nREACTION: %+v", reaction)
	Observe(action, reaction, perspective)
	return reaction
}

func main() {
	fmt.Printf("\n%+v", World("hi"))
}
