package game

import (
	"fmt"
	"regexp"
	"strings"
)

func HangmanAgreement(world World, action Action) (r Reaction) {
	actionLength := len(action)
	worldString := string(world)
	actionString := string(action)

	if actionLength > 1 && actionString == worldString {
		return Reaction{Response: Response(world), Reward: 1.0}

	} else if actionLength > 1 && actionString != worldString {
		return Reaction{Response: Response("!"), Reward: 0.0}

	} else if strings.Contains(worldString, actionString) {
		re := regexp.MustCompile(fmt.Sprintf("[^%s]", actionString))
		response := re.ReplaceAllString(worldString, "-")
		return Reaction{Response: Response(response), Reward: 1.0}

	} else {
		return Reaction{Response: Response("!"), Reward: 0.0}

	}
}
