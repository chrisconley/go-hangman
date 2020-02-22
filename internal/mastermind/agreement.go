package mastermind

import (
	// "fmt"
	"regexp"
	"strings"

	"github.com/chrisconley/go-hangman/internal/game"
)
func Agreement(world game.World, action game.Action) (r game.Reaction) {
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

	worldString := string(world)
	actionString := string(action)

	response := ""
	reward := 0.0

	for i := 0; i < len(actionString); i++ {
		if worldString[i] == actionString[i] {
			response = response + "B"
			reward += 0.25
			worldString = worldString[:i] + "-" + worldString[i+1:]
			actionString = actionString[:i] + "-" + actionString[i+1:]
		}
	}

	for i := 0; i < len(actionString); i++ {
		if actionString[i] == '-' {
			continue
		}

		re := regexp.MustCompile(string(actionString[i]))
		match := re.FindString(worldString)
		if match != "" {
			response = response + "W"
			reward += 0.125
			worldString = strings.Replace(worldString, match, "-", 1)
		}
	}

	return game.Reaction{Action: action, Response: game.Response(response), Reward: reward}
}