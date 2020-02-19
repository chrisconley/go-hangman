package game

import "strings"

func TestAgreement(w World, a Action) (r Reaction) {
	if strings.Contains(string(w), string(a)) {
		return Reaction{Action: a, Response: Response(a), Reward: 1.0}
	} else {
		return Reaction{Action: a, Response: "!", Reward: 0.0}
	}
}
