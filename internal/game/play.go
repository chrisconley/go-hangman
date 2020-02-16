package game

import (
	"fmt"
	// "math"
	// "regexp"
)

// ### The system ###
type World string
type Action string
type Response string
type Reaction struct {
	Action Action
	Response Response
	Reward   float64
}
type Agreement func(w World, action Action) (reaction Reaction)

type Knowledge []Reaction

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
	Knowledge       Knowledge
	PossibleWorlds  PossibleWorlds
	PossibleActions PossibleActions
	Agreement		Agreement
	Motivation      Motivation
}

func (p *Perspective) UpdateKnowledge(reaction Reaction) {
	p.Knowledge = append(p.Knowledge, reaction)
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
	fmt.Printf("\nPOSSIBLE ACTIONS: %+v\n", perspective.PossibleActions)
	fmt.Printf("\nPOSSIBLE WORLDS: %+v\n", perspective.PossibleWorlds)

}

// ### Process implementation ###
// func PlaceholderHangmanAgreement(w World, d Action) (r Reaction) { return Reaction{"hi", 1.0} }

func PlayGame(truth World, perspective *Perspective) (World) {
	for true {
		playTurn(truth, perspective)
		if len(perspective.PossibleWorlds) == 1 {
			break
		}
	}
	fmt.Printf("\nGAME OVER: %+v\n", perspective.PossibleWorlds[0])
	return perspective.PossibleWorlds[0].World

}

func playTurn(truth World, perspective *Perspective) (reaction Reaction) {
	predictions := Explore(*perspective, HangmanAgreement)
	action := Decide(predictions, perspective.Motivation)
	fmt.Printf("\nACTION: %+v", action)
	reaction = Act(action, HangmanAgreement, truth)
	fmt.Printf("\nREACTION: %+v", reaction)
	Observe(action, reaction, perspective)
	return reaction
}

func main() {
	truth := World("abc")
	perspective := &Perspective{
		Knowledge:       Knowledge{},
		PossibleWorlds:  PossibleWorlds{PossibleWorld{"abc", 0.5}, PossibleWorld{"cab", 0.5}},
		PossibleActions: PossibleActions{"a", "b", "c"},
		Motivation:      Motivation{0.5, 0.25, 0.25},
	}
	// Where does the Goal come into play? Is that implicit in the Agreement?
	predictions := Explore(*perspective, HangmanAgreement)
	action := Decide(predictions, perspective.Motivation)
	reaction := Act(action, HangmanAgreement, truth)
	Observe(action, reaction, perspective)
	fmt.Printf("\n%+v", perspective.Knowledge)
}
