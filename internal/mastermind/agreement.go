package mastermind

import (
	"bytes"
	"sort"
	"strings"

	"github.com/chrisconley/go-hangman/internal/game"
)

// GetActions offers every code not played yet, except codes that are certain to split the remaining
// codes exactly as an offered code does.
//
// Two codes split the remaining codes the same way when one turns into the other by moving pegs to
// other positions and renaming colours, and that move leaves every guess played so far unchanged and
// turns the list of codes into itself. Such a move turns each remaining code into another remaining code
// and keeps every response the same, so the two codes give the same groups, the same scores, and either
// both could be the answer or neither could. When the list holds every combination of the colours,
// colours that appear in no guess played so far can always be renamed among themselves. Pegs can be
// moved only when every played guess, with its colours renamed, looks the same afterwards, as before the
// first guess.
//
// Of each set of codes that split the same way, the alphabetically first is offered. That is the code
// Decide picks when the whole set ties, so play is the same as when every code is offered.
func GetActions(gameplay game.GamePlay, remainingWorlds game.PossibleWorlds, aa game.AvailableActions) []game.Action {
	codes := aa.WorldActions()
	if len(codes) == 0 {
		return []game.Action{}
	}
	length := codes[0].Length()

	played := make([]string, len(gameplay.Reactions))
	isPlayedColor := [256]bool{}
	for i, reaction := range gameplay.Reactions {
		played[i] = reaction.Action.String()
		for _, b := range reaction.Action.Bytes() {
			isPlayedColor[b] = true
		}
	}
	unplayedColors := []byte{}
	for _, color := range aa.Actions() {
		if !isPlayedColor[color.Bytes()[0]] {
			unplayedColors = append(unplayedColors, color.Bytes()[0])
		}
	}
	sort.Slice(unplayedColors, func(i, j int) bool { return unplayedColors[i] < unplayedColors[j] })

	moves := movesKeepingGuesses(played, length)
	freelyRenamed := unplayedColors
	if len(codes) != intPower(len(aa.Actions()), length) {
		moves = movesKeepingCodes(moves, unplayedColors, codes)
		freelyRenamed = nil
	}
	isFree := [256]bool{}
	for _, color := range freelyRenamed {
		isFree[color] = true
	}

	actions := []game.Action{}
	var another game.Action
	moved := make([]byte, length)
	for _, code := range codes {
		if gameplay.PreviouslyActed(code) {
			continue
		}
		if isFirstOfSameSplit(code.Bytes(), moves, isFree, freelyRenamed, moved) {
			actions = append(actions, code)
		} else if another == nil {
			another = code
		}
	}
	// Decide needs two guesses to choose between. When every code splits the same way, a second one
	// from the same set is offered; it scores exactly the same and comes later alphabetically, so it is
	// never picked.
	if len(actions) == 1 && another != nil {
		actions = append(actions, another)
	}
	return actions
}

// move puts the peg at position p in position to[p] and renames each colour c to rename[c].
type move struct {
	to     []int
	rename [256]byte
}

func unchangedColors() [256]byte {
	rename := [256]byte{}
	for c := range rename {
		rename[c] = byte(c)
	}
	return rename
}

// mostPositionOrders limits how many orders of the pegs are tried: 720 is every order of 6 pegs. Longer
// codes have too many orders to try, so for them pegs stay where they are and only colours are renamed.
const mostPositionOrders = 720

// movesKeepingGuesses lists every way of moving pegs to other positions, with a renaming of the played
// colours, that leaves every played guess unchanged. Leaving the pegs where they are is always one.
func movesKeepingGuesses(played []string, length int) []move {
	unmoved := make([]int, length)
	for p := range unmoved {
		unmoved[p] = p
	}
	orders := [][]int{unmoved}
	if cappedFactorial(length, mostPositionOrders+1) <= mostPositionOrders {
		orders = positionOrders(length)
	}
	moves := []move{}
	for _, to := range orders {
		rename := unchangedColors()
		renamedFrom := unchangedColors()
		isSet := [256]bool{}
		isTaken := [256]bool{}
		keeps := true
		for _, guess := range played {
			for p := 0; p < length && keeps; p++ {
				from, into := guess[p], guess[to[p]]
				if (isSet[from] && rename[from] != into) || (isTaken[into] && renamedFrom[into] != from) {
					keeps = false
				}
				rename[from], isSet[from] = into, true
				renamedFrom[into], isTaken[into] = from, true
			}
		}
		if keeps {
			moves = append(moves, move{to: to, rename: rename})
		}
	}
	return moves
}

// movesKeepingCodes is for a list of codes that does not hold every combination of the colours. It
// adds every renaming of the unplayed colours to each move, and keeps the moves that turn the list into
// itself. When that would take too long to check, it keeps only the move that changes nothing, so every
// code is offered.
func movesKeepingCodes(moves []move, unplayedColors []byte, codes []game.Action) []move {
	const mostSteps = 10_000_000
	if len(moves)*len(codes)*cappedFactorial(len(unplayedColors), mostSteps+1) > mostSteps {
		return []move{moves[0]}
	}
	renamings := colorOrders(unplayedColors)
	isCode := map[string]bool{}
	for _, code := range codes {
		isCode[code.String()] = true
	}
	kept := []move{}
	moved := make([]byte, codes[0].Length())
	for _, m := range moves {
		for _, order := range renamings {
			full := move{to: m.to, rename: m.rename}
			for i, color := range unplayedColors {
				full.rename[color] = order[i]
			}
			keepsCodes := true
			for _, code := range codes {
				applyMove(code.Bytes(), full, moved)
				if !isCode[string(moved)] {
					keepsCodes = false
					break
				}
			}
			if keepsCodes {
				kept = append(kept, full)
			}
		}
	}
	return kept
}

func applyMove(code []byte, m move, moved []byte) {
	for p, color := range code {
		moved[m.to[p]] = m.rename[color]
	}
}

// isFirstOfSameSplit tells whether code comes before, alphabetically, every code that the moves,
// followed by any renaming of the colours in freelyRenamed among themselves, make out of it. moved is
// space for the work.
func isFirstOfSameSplit(code []byte, moves []move, isFree [256]bool, freelyRenamed []byte, moved []byte) bool {
	for _, m := range moves {
		applyMove(code, m, moved)
		renamed := [256]byte{}
		next := 0
		for p, color := range moved {
			if !isFree[color] {
				continue
			}
			if renamed[color] == 0 {
				renamed[color] = freelyRenamed[next]
				next++
			}
			moved[p] = renamed[color]
		}
		if bytes.Compare(moved, code) < 0 {
			return false
		}
	}
	return true
}

// cappedFactorial returns n! when it is below most, and most otherwise.
func cappedFactorial(n int, most int) int {
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
		if result >= most {
			return most
		}
	}
	return result
}

// positionOrders lists every order of the positions 0 to length-1, starting with 0, 1, 2, ...
func positionOrders(length int) [][]int {
	positions := make([]int, length)
	for i := range positions {
		positions[i] = i
	}
	return orders(positions)
}

func colorOrders(colors []byte) [][]byte {
	asInts := make([]int, len(colors))
	for i, c := range colors {
		asInts[i] = int(c)
	}
	result := [][]byte{}
	for _, order := range orders(asInts) {
		colorOrder := make([]byte, len(order))
		for i, c := range order {
			colorOrder[i] = byte(c)
		}
		result = append(result, colorOrder)
	}
	return result
}

// orders lists every order of items, starting with items as given.
func orders(items []int) [][]int {
	if len(items) <= 1 {
		return [][]int{append([]int{}, items...)}
	}
	result := [][]int{}
	for i := range items {
		rest := append(append([]int{}, items[:i]...), items[i+1:]...)
		for _, order := range orders(rest) {
			result = append(result, append([]int{items[i]}, order...))
		}
	}
	return result
}

func intPower(base int, exponent int) int {
	result := 1
	for i := 0; i < exponent; i++ {
		result *= base
	}
	return result
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
