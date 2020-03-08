package hangman

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chrisconley/go-hangman/internal/game"
)

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
