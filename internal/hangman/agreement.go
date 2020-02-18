package hangman

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chrisconley/go-hangman/internal/game"
)

func UpdatePossibleActions(reaction game.Reaction, possibleActions game.PossibleActions) (a game.PossibleActions) {
	a = game.PossibleActions{}
	for _, possibleAction := range possibleActions {
		if possibleAction != reaction.Action {
			a = append(a, possibleAction)
		}
	}
	return a
}

func UpdatePossibleWorlds(reaction game.Reaction, possibleWorlds game.PossibleWorlds) (p game.PossibleWorlds) {
	p = game.PossibleWorlds{}
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

func HangmanAgreement(world game.World, action game.Action) (r game.Reaction) {
	// fmt.Printf("\nACTION IN HANGMENTAGREEMENT: %+v", action)
	actionLength := len(action)
	worldString := string(world)
	actionString := string(action)

	if actionLength > 1 && actionString == worldString {
		return game.Reaction{Action: action, Response: game.Response(world), Reward: 1.0}

	} else if actionLength > 1 && actionString != worldString {
		return game.Reaction{Action: action, Response: game.Response("!"), Reward: 0.0}

	} else if strings.Contains(worldString, actionString) {
		re := regexp.MustCompile(fmt.Sprintf("[^%s]", actionString))
		response := re.ReplaceAllString(worldString, "-")
		return game.Reaction{Action: action, Response: game.Response(response), Reward: 1.0}

	} else {
		return game.Reaction{Action: action, Response: game.Response("!"), Reward: 0.0}

	}
}
