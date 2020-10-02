package hangman

import (
	"strings"

	"github.com/chrisconley/go-hangman/internal/game"
)

// TODO: Make this take game.PossibleWorlds
func GetActions(remainingWorlds game.PossibleWorlds, aa game.AvailableActions) []game.Action {
	// Only use world actions when there are only a few choices left
	if remainingWorlds.Len() < aa.ActionsLength() {
		worldActions := []game.Action{}
		for _, index := range remainingWorlds.Indices() {
			world := game.BUniverse.Worlds[index]
			worldActions = append(worldActions, game.Action(world.Name))
		}
		return append(aa.Actions(), worldActions...)
	} else {
		return aa.Actions()
	}
}

func Act(world game.World, action game.Action) (r game.Reaction) {
	actionLength := len(action)
	worldString := string(world)
	actionString := string(action)

	if actionLength > 1 && actionString == worldString {
		return game.Reaction{Action: action, Response: game.Response(world), Reward: 1.0}
	} else if actionLength > 1 && actionString != worldString {
		return game.Reaction{Action: action, Response: game.Response("!"), Reward: 0.0}
	} else {
		hasMatch := false
		var buffer strings.Builder
		for _, character := range worldString {
			if string(character) == actionString {
				buffer.WriteString(actionString)
				hasMatch = true
			} else {
				buffer.WriteString("-")
			}
		}
		if hasMatch {
			return game.Reaction{Action: action, Response: game.Response(buffer.String()), Reward: 1.0}
		} else {
			return game.Reaction{Action: action, Response: game.Response("!"), Reward: 0.0}
		}

	}
}

var Agreement = game.Agreement{
	GetActions: GetActions,
	Act:        Act,
}
