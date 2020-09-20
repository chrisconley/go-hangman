package game

import "strings"

func BuildTestAgreement(reward float64) func(World, Action) Reaction {
	return func(w World, a Action) (r Reaction) {
		if strings.Contains(string(w), string(a)) {
			return Reaction{Action: a, Response: Response(a), Reward: reward}
		} else {
			return Reaction{Action: a, Response: "!", Reward: 0.0}
		}
	}

}

func TestActionFilter(remainingWorlds PossibleWorlds, aa AvailableActions) []Action {
	return []Action{"hi"}
}
