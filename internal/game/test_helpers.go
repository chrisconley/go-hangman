package game

import "strings"

func BuildTestAct(reward float64) func(World, Action) Reaction {
	return func(w World, a Action) (r Reaction) {
		if strings.Contains(string(w), string(a)) {
			return Reaction{Action: a, Response: Response(a), Reward: reward}
		} else {
			return Reaction{Action: a, Response: "!", Reward: 0.0}
		}
	}
}
