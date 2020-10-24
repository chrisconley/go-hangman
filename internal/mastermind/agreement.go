package mastermind

import (
	"sort"
	"strings"

	"github.com/chrisconley/go-hangman/internal/game"
)

type kv struct {
	Key   string
	Value int
}

func GetActions(gameplay game.GamePlay, remainingWorlds game.PossibleWorlds, aa game.AvailableActions) []game.Action {
	allColors := make([]string, len(aa.Actions()))
	for index, action := range aa.Actions() {
		allColors[index] = action.String()
	}
	usedColors := make(map[string]int8)
	for _, reaction := range gameplay.Reactions {
		for _, b := range reaction.Action.Bytes() {
			usedColors[string(b)] = 1

		}
	}
	freeColors := make(map[string]int8)
	freeColorsSlice := make([]string, (len(allColors) - len(usedColors)))
	freeColorIndex := 0
	for _, color := range allColors {
		if _, ok := usedColors[color]; !ok {
			freeColors[color] = 1
			freeColorsSlice[freeColorIndex] = color
			freeColorIndex++
		}
	}

	builtActionStrings := make(map[string]int8)
	for _, worldAction := range aa.WorldActions() {
		symActionColorTypes := make([]string, worldAction.Length())
		freeColorCounts := make(map[string]int)

		for index, b := range worldAction.Bytes() {
			char := string(b)
			if _, ok := freeColors[char]; !ok {
				symActionColorTypes[index] = "c"
			} else {
				freeColorCounts[char]++
				symActionColorTypes[index] = "f"
			}

		}
		sortedCounts := []kv{}
		for k, v := range freeColorCounts {
			sortedCounts = append(sortedCounts, kv{k, v})
		}

		sort.Slice(sortedCounts, func(i, j int) bool {
			return sortedCounts[i].Value > sortedCounts[j].Value
		})

		freeColorsToAssign := []string{}
		for index, kv := range sortedCounts {
			for j := 0; j < kv.Value; j++ {
				freeColorsToAssign = append(freeColorsToAssign, freeColorsSlice[index])
			}
		}

		builtActionBytes := make([]byte, worldAction.Length())
		freeColorIndex := 0
		for index, colorType := range symActionColorTypes {
			if colorType == "c" {
				builtActionBytes[index] = worldAction.Bytes()[index]
			} else if colorType == "f" {
				builtActionBytes[index] = freeColorsToAssign[freeColorIndex][0]
				freeColorIndex++
			}
		}

		builtActionStrings[string(builtActionBytes)] = 1
	}
	actions := []game.Action{}
	for actionString, _ := range builtActionStrings {
		action := game.NewAction(actionString)
		if !gameplay.PreviouslyActed(action) {
			actions = append(actions, action)
		}

	}
	return actions
}

func Act(gameplay game.GamePlay, world game.World, action game.Action) (r game.Reaction) {
	// The mastermind rules can be a bit ambiguous.

	// The code below is implemented so that duplicate white responses are not returned.
	// This matches the rules of the game and what Knuth implemented his original minimax solution
	// and subsequent papers in resources/ use:
	// 6 garners a "white" response
	// 3 garners a "white" response
	// 2 garners a "white" response
	// 6 does not garner a "white" response becuase the "6" in the code word was already "used"
	//     for the first "6" in the guess
	// This results in a response of 'WWW'.

	if gameplay.PreviouslyActed(action) {
		panic("This action is not allowed because it has already been used")
	}

	worldString := world.String()
	actionString := action.String()

	responseBuffer := strings.Builder{}
	reward := 0.0
	totalPossibleReward := float64(len(worldString))

	remainingActionCharacters := strings.Builder{}
	remainingWorldCharacters := make(map[byte]int)

	// First, we're going to look for "B" responses (guesses in the correction position)
	for i := range actionString {
		if worldString[i] == actionString[i] {
			// If the character in position i is the same in both
			// Add a "B" to the response and record reward
			responseBuffer.WriteString("B")
			reward += 1 / totalPossibleReward
		} else {
			// Otherwise, save the character to check for "W" matches
			// and increase the count of the world character
			remainingActionCharacters.WriteByte(actionString[i])
			remainingWorldCharacters[worldString[i]]++
		}
	}

	// Now iterate through the action characters that were not already used for a "B" match
	for _, char := range remainingActionCharacters.String() {
		character := byte(char)
		// Look up the action character to see if there are any matching
		// world characters left
		if remainingWorldCharacters[character] > 0 {
			// If there are, record the "W" response and reward
			responseBuffer.WriteString("W")
			reward += 0.5 / totalPossibleReward
			// Reduce the count to record that that world character has been used
			// so we don't double count it in a future iteration
			remainingWorldCharacters[character]--

		}
	}

	response := responseBuffer.String()
	if response == "" {
		response = "!"
	}
	reaction := game.NewReaction(action, response, reward)
	return reaction
}

var Agreement = game.Agreement{
	GetActions: GetActions,
	Act:        Act,
}
