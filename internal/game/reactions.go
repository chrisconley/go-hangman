package game

import (
	"github.com/chrisconley/go-hangman/internal/utils"
)

type reactionExplorer struct {
	tracker map[Action]actionData
	actions []Action
}

type actionData struct {
	Responses map[Response]*PossibleWorlds
	Counts    *utils.GroupedIncrementer
	Rewards   *utils.GroupedIncrementer
}

func NewReactionExplorer() reactionExplorer {
	re := reactionExplorer{}
	re.tracker = make(map[Action]actionData)
	re.actions = make([]Action, 0)
	return re
}

func (re *reactionExplorer) Add(action Action, response Response, reward float64, world World) {
	actionData, responsePossibleWorlds := re.findOrInitializeResponse(action, response)
	actionData.Counts.Increment(string(response), 1.0)
	actionData.Rewards.Increment(string(response), reward)
	// fmt.Printf("possible worlds pointer %p\n", &responsePossibleWorlds)
	// fmt.Printf("tracker pointer %p\n", &re.tracker)
	responsePossibleWorlds.Add(world)
}

func (re *reactionExplorer) Get(action Action, response Response) PossibleWorlds {
	_, responsePossibleWorlds := re.findOrInitializeResponse(action, response)
	return *responsePossibleWorlds
}

func (re *reactionExplorer) ActionCounts(action Action) []float64 {
	actionData := re.findOrInitializeAction(action)
	return actionData.Counts.All()
}

func (re *reactionExplorer) ActionTotal(action Action) float64 {
	actionData := re.findOrInitializeAction(action)
	return actionData.Counts.Total()
}

func (re *reactionExplorer) ActionRewardTotal(action Action) float64 {
	actionData := re.findOrInitializeAction(action)
	return actionData.Rewards.Total()
}

func (re *reactionExplorer) Actions() []Action {
	return re.actions
}

func (re *reactionExplorer) findOrInitializeAction(action Action) actionData {
	_, ok := re.tracker[action]
	if !ok {
		re.actions = append(re.actions, action)
		counts := utils.NewGroupedIncrementer()
		rewards := utils.NewGroupedIncrementer()
		re.tracker[action] = actionData{
			Responses: make(map[Response]*PossibleWorlds),
			Counts:    &counts,
			Rewards:   &rewards,
		}
	}
	return re.tracker[action]
}

func (re *reactionExplorer) findOrInitializeResponse(action Action, response Response) (actionData, *PossibleWorlds) {
	actionData := re.findOrInitializeAction(action)
	_, ok := re.tracker[action].Responses[response]
	if !ok {
		possibleWorlds := NewPossibleWorlds()
		re.tracker[action].Responses[response] = &possibleWorlds
	}
	return actionData, actionData.Responses[response]
}
