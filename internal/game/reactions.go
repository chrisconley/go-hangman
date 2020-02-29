package game

import (
	"fmt"

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
	Counts    utils.GroupedIncrementer
	Rewards   utils.GroupedIncrementer
}

type ReactionExplorer struct {
	tracker map[Action]actionData
}

func NewReactionExplorer() ReactionExplorer {
	re := ReactionExplorer{}
	re.tracker = make(map[Action]actionData)
	return re
}

func (re *ReactionExplorer) Add(action Action, response Response, rw float64, world World) {
	re.findOrInitializeResponse(action, response)
	re.tracker[action].Responses[response] = append(re.tracker[action].Responses[response], world)
}

func (re *ReactionExplorer) Get(action Action, response Response) []World {
	re.findOrInitializeResponse(action, response)
	return re.tracker[action].Responses[response]
}

func (re *ReactionExplorer) ActionCounts(a Action) []float64 {
	return []float64{1.1117}
}

func (re *ReactionExplorer) ActionTotal(a Action) float64 {
	return 1.1117
}

func (re *ReactionExplorer) ActionRewardTotal(a Action) float64 {
	return 1.1117
}

func (re *ReactionExplorer) findOrInitializeAction(action Action) {
	_, ok := re.tracker[action]
	if !ok {
		re.tracker[action] = actionData{Responses: make(map[Response][]World)}
	}
}

func (re *ReactionExplorer) findOrInitializeResponse(action Action, response Response) {
	re.findOrInitializeAction(action)
	_, ok := re.tracker[action].Responses[response]
	if !ok {
		fmt.Printf("hi")
		re.tracker[action].Responses[response] = make([]World, 0)
	}
}
