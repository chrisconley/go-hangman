package hangman

import (
	"strings"

	"github.com/chrisconley/go-hangman/internal/game"
)

func GetActions(gameplay game.GamePlay, remainingWorlds game.PossibleWorlds, aa game.AvailableActions) []game.Action {
	actions := []game.Action{}
	for _, action := range aa.Actions() {
		if !gameplay.PreviouslyActed(action) {
			actions = append(actions, action)
		}
	}

	// Only use world actions when there are only a few choices left
	if remainingWorlds.Len() < aa.ActionsLength() {
		worldActions := []game.Action{}
		for _, index := range remainingWorlds.Indices() {
			world := game.BUniverse.GetWorldByIndex(index)
			worldAction := game.Action(world.Name)
			if !gameplay.PreviouslyActed(worldAction) {
				worldActions = append(worldActions, worldAction)
			}
		}
		return append(actions, worldActions...)
	} else {
		return actions
	}
}

func Act(gameplay game.GamePlay, world game.World, action game.Action) (r game.Reaction) {
	if gameplay.PreviouslyActed(action) {
		panic("This action is not allowed because it has already been used")
	}

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
