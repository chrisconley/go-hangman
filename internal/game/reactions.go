package game

type ResponseCounts map[Response]*Counts

type Counts struct {
	Weight float32
	Count  int
}
type Reward struct {
	Utility float64
	Regret  float64
}
type ResponseCounter struct {
	ResponseCounts ResponseCounts
	Reward         Reward
}
type ReactionCounter map[Action]ResponseCounter

func (p ResponseCounter) CountTotal() int {
	total := int(0)
	for _, count := range p.ResponseCounts {
		total += count.Count
	}
	return total
}

func (p ResponseCounter) AllCounts() []int {
	counts := make([]int, 0, len(p.ResponseCounts))
	for _, count := range p.ResponseCounts {
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
		rc[action] = ResponseCounter{ResponseCounts: ResponseCounts{}, Reward: Reward{0, 0}}
	}
}

func (rc ReactionCounter) findOrInitializeResponse(action Action, response Response) {
	rc.findOrInitializeAction(action)
	_, ok := rc[action].ResponseCounts[response]
	if ok {

	} else {
		rc[action].ResponseCounts[response] = &Counts{Weight: 0, Count: 0}
	}
}

func (rc ReactionCounter) Get(r Reaction) int {
	rc.findOrInitializeResponse(r.Action, r.Response)
	return rc[r.Action].ResponseCounts[r.Response].Count
}

func (rc ReactionCounter) Increment(r Reaction) {
	rc.findOrInitializeResponse(r.Action, r.Response)
	rc[r.Action].ResponseCounts[r.Response].Count += 1
}

func (rc ReactionCounter) Set(r Reaction, weight float32, count int) {
	rc.findOrInitializeResponse(r.Action, r.Response)
	rc[r.Action].ResponseCounts[r.Response] = &Counts{Weight: weight, Count: count}
}
