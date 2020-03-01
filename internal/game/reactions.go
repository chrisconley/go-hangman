package game

import (
	"github.com/chrisconley/go-hangman/internal/utils"
)

type ResponseCount struct {
	Count  int
	Reward float64
}

type ResponseCounter map[Response]*ResponseCount

type ReactionCounter map[Action]ResponseCounter

func (p ResponseCounter) RewardTotal() float64 {
	total := 0.0
	for _, count := range p {
		total += count.Reward
	}
	return total
}

func (p ResponseCounter) CountTotal() int {
	total := int(0)
	for _, count := range p {
		total += count.Count
	}
	return total
}

func (p ResponseCounter) AllCounts() []int {
	counts := make([]int, 0, len(p))
	for _, count := range p {
		counts = append(counts, count.Count)
	}
	return counts
}

func (rc ReactionCounter) Actions() []Action {
	actions := []Action{}
	for action, _ := range rc {
		actions = append(actions, action)
	}
	return actions
}

func (rc ReactionCounter) TotalRewardForAction(action Action) float64 {
	rc.findOrInitializeAction(action)
	return rc[action].RewardTotal()
}

func (rc ReactionCounter) TotalCountForAction(action Action) int {
	rc.findOrInitializeAction(action)
	return rc[action].CountTotal()
}

func (rc ReactionCounter) AllCountsForAction(action Action) []int {
	rc.findOrInitializeAction(action)
	return rc[action].AllCounts()
}
func (rc ReactionCounter) findOrInitializeAction(action Action) {
	_, ok := rc[action]
	if ok {

	} else {
		rc[action] = ResponseCounter{}
	}
}

func (rc ReactionCounter) findOrInitializeResponse(action Action, response Response) {
	rc.findOrInitializeAction(action)
	_, ok := rc[action][response]
	if ok {

	} else {
		rc[action][response] = &ResponseCount{Count: 0, Reward: 0.0}
	}
}

func (rc ReactionCounter) Get(r Reaction) int {
	rc.findOrInitializeResponse(r.Action, r.Response)
	return rc[r.Action][r.Response].Count
}

func (rc ReactionCounter) Increment(r Reaction) {
	rc.findOrInitializeResponse(r.Action, r.Response)
	rc[r.Action][r.Response].Count += 1
	rc[r.Action][r.Response].Reward += r.Reward
}

func (rc ReactionCounter) Set(r Reaction, reward float64, count int) {
	rc.findOrInitializeResponse(r.Action, r.Response)
	rc[r.Action][r.Response] = &ResponseCount{Count: count, Reward: reward}
}

type actionData struct {
	Responses map[Response][]World
	Counts    *utils.GroupedIncrementer
	Rewards   *utils.GroupedIncrementer
}

type reactionExplorer struct {
	tracker map[Action]actionData
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
