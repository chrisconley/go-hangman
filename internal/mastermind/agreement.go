package mastermind

import (
	"strings"

	"github.com/chrisconley/go-hangman/internal/game"
)

func GetActions(gameplay game.GamePlay, remainingWorlds game.PossibleWorlds, aa game.AvailableActions) []game.Action {
	actions := []game.Action{}
	for _, action := range aa.WorldActions() {
		if !gameplay.PreviouslyActed(action) {
			actions = append(actions, action)
		}
	}
	return actions
}

func Act(gameplay game.GamePlay, world game.World, action game.Action) (r game.Reaction) {
	// In An Optimal Mastermind (4,7) Strategy and More Results in the Expected Case
	// https://arxiv.org/pdf/1305.1010.pdf, the author "double" counts duplicate guesses
	// that garner a "white" response.

	// Example: The code word is 3632. Let's work out the response for guess 6326:
	// 6 garners a "white" response
	// 3 garners a "white" response
	// 2 garners a "white" response
	// 6 garners a "white" response again
	// This results in a response of 'WWWW'.

	// In contrast, the code below is implemented so that duplicate white responses are not returned.
	// This matches the rules of the game and what Knuth implemented his original minimax solution:
	// 6 garners a "white" response
	// 3 garners a "white" response
	// 2 garners a "white" response
	// 6 does not garner a "white" response becuase the "6" in the code word was already "used"
	//     for the first "6" in the guess
	// This results in a response of 'WWW'.

	if gameplay.PreviouslyActed(action) {
		panic("This action is not allowed because it has already been used")
	}

	worldString := string(world)
	actionString := string(action)
	worldLength := len(worldString)

	var responseBuffer strings.Builder
	reward := 0.0

	var remainingActions strings.Builder
	remainingWorldCharacters := make(map[byte]int)

	// First, we're going to look for "B" responses (guesses in the correction position)
	for i := range actionString {
		if worldString[i] == actionString[i] {
			responseBuffer.WriteString("B")
			reward += 1 / float64(worldLength)
		} else {
			remainingActions.WriteByte(actionString[i])
			remainingWorldCharacters[worldString[i]]++
		}
	}

	for _, char := range remainingActions.String() {
		character := byte(char)
		if remainingWorldCharacters[character] > 0 {
			responseBuffer.WriteString("W")
			reward += 0.5 / float64(worldLength)
			remainingWorldCharacters[character]--

		}
	}

	response := responseBuffer.String()
	if response == "" {
		response = "!"
	}
	return game.Reaction{Action: action, Response: game.Response(response), Reward: reward}
}

var Agreement = game.Agreement{
	GetActions: GetActions,
	Act:        Act,
}
