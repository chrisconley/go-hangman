package game

import (
	"fmt"
	"regexp"
	"strings"
)

func UpdatePossibleActions(reaction Reaction, possibleActions PossibleActions) (a PossibleActions) {
	a = PossibleActions{}
	for _, possibleAction := range possibleActions {
		if possibleAction != reaction.Action {
			a = append(a, possibleAction)
		}
	}
	return a
}

func UpdatePossibleWorlds(reaction Reaction, possibleWorlds PossibleWorlds) (p PossibleWorlds) {
	p = PossibleWorlds{}
	if reaction.Response == "!" {
		re := regexp.MustCompile(string(reaction.Action))
		for _, possibleWorld := range possibleWorlds {
			worldString := string(possibleWorld.World)
			if !re.MatchString(worldString) {
				p = append(p, possibleWorld)
			}
		}
	} else {
		re := regexp.MustCompile("-")
		pattern := re.ReplaceAllString(string(reaction.Response), ".")
		re = regexp.MustCompile(pattern)

		for _, possibleWorld := range possibleWorlds {
			worldString := string(possibleWorld.World)
			if re.MatchString(worldString) {
				p = append(p, possibleWorld)
			}

		}
	}
	return p
}

func HangmanAgreement(world World, action Action) (r Reaction) {
	// fmt.Printf("\nACTION IN HANGMENTAGREEMENT: %+v", action)
	actionLength := len(action)
	worldString := string(world)
	actionString := string(action)

	if actionLength > 1 && actionString == worldString {
		return Reaction{Action: action, Response: Response(world), Reward: 1.0}

	} else if actionLength > 1 && actionString != worldString {
		return Reaction{Action: action, Response: Response("!"), Reward: 0.0}

	} else if strings.Contains(worldString, actionString) {
		re := regexp.MustCompile(fmt.Sprintf("[^%s]", actionString))
		response := re.ReplaceAllString(worldString, "-")
		return Reaction{Action: action, Response: Response(response), Reward: 1.0}

	} else {
		return Reaction{Action: action, Response: Response("!"), Reward: 0.0}

	}
}
