package hangman

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chrisconley/go-hangman/internal/game"
)

// TODO: Make sure this results in optimal number of turns appropriately
func ActionFilter(remainingWorlds game.Worlds, aa game.AvailableActions) []game.Action {
	// Only use world actions when there are only a few choices left
	if len(remainingWorlds) < aa.ActionsLength() {
		worldActions := []game.Action{}
		for _, world := range remainingWorlds {
			worldActions = append(worldActions, game.Action(world))
		}
		return append(aa.Actions(), worldActions...)
	} else {
		return aa.Actions()
	}

	// Use all actions
	// return aa.AllActions()

	// Use just letter action
	// return aa.Actions()
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
