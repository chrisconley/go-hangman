package game

import "strings"

func BuildActions(strings ...string) []Action {
	actions := make([]Action, len(strings))
	for index, actionString := range strings {
		actions[index] = NewAction(actionString)
	}
	return actions
}

func BuildTestAct(reward float64) func(GamePlay, World, Action) Reaction {
	return func(gp GamePlay, w World, a Action) (r Reaction) {
		if strings.Contains(string(w), a.String()) {
			return Reaction{Action: a, Response: Response(a.String()), Reward: reward}
		} else {
			return Reaction{Action: a, Response: "!", Reward: 0.0}
		}
	}
}
