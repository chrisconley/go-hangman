package game

import (
	"github.com/chrisconley/go-hangman/internal/utils"
)

type reactionExplorer struct {
	tracker map[Action]actionData
}

type actionData struct {
	Responses map[Response][]World
	Counts    *utils.GroupedIncrementer
	Rewards   *utils.GroupedIncrementer
}

func NewReactionExplorer() reactionExplorer {
	re := reactionExplorer{}
	re.tracker = make(map[Action]actionData)
	return re
}

func (re *reactionExplorer) Add(action Action, response Response, reward float64, world World) {
	re.findOrInitializeResponse(action, response)
	counts := re.tracker[action].Counts
	counts.Increment(string(response), 1.0)
	rewards := re.tracker[action].Rewards
	rewards.Increment(string(response), reward)
	re.tracker[action].Responses[response] = append(re.tracker[action].Responses[response], world)
}

func (re *reactionExplorer) Get(action Action, response Response) []World {
	re.findOrInitializeResponse(action, response)
	return re.tracker[action].Responses[response]
}

func (re *reactionExplorer) ActionCounts(action Action) []float64 {
	re.findOrInitializeAction(action)
	return re.tracker[action].Counts.All()
}

func (re *reactionExplorer) ActionTotal(action Action) float64 {
	re.findOrInitializeAction(action)
	return re.tracker[action].Counts.Total()
}

func (re *reactionExplorer) ActionRewardTotal(action Action) float64 {
	re.findOrInitializeAction(action)
	return re.tracker[action].Rewards.Total()
}

func (re *reactionExplorer) Actions() []Action {
	actions := []Action{}
	for action, _ := range re.tracker {
		actions = append(actions, action)
	}
	return actions
}

func (re *reactionExplorer) findOrInitializeAction(action Action) {
	_, ok := re.tracker[action]
	if !ok {
		counts := utils.NewGroupedIncrementer()
		rewards := utils.NewGroupedIncrementer()
		re.tracker[action] = actionData{
			Responses: make(map[Response][]World),
			Counts:    &counts,
			Rewards:   &rewards,
		}
	}
}

func (re *reactionExplorer) findOrInitializeResponse(action Action, response Response) {
	re.findOrInitializeAction(action)
	_, ok := re.tracker[action].Responses[response]
	if !ok {
		re.tracker[action].Responses[response] = make([]World, 0)
	}
}
