// Package wordle supplies the rules of Wordle to the engine. It plugs in the same way the hangman and
// mastermind packages do, with a game.Agreement: the pair of rules saying which guesses may be played and
// how a guess is answered. The engine's turn loop (score every guess, pick one, see the answer to it,
// cross off the answers that no longer fit) then plays Wordle without changes. The possible answers come
// one per line on stdin, as for the other games. The allowed guesses are passed to NewAgreement, because
// the engine otherwise only knows the possible answers. Every answer must also be an allowed guess, and
// every word must be five lowercase letters.
//
// The response to a guess is one character per position: G when the letter is in the answer at that
// position, Y when it is in the answer somewhere else, and - when it is not. A letter that appears more
// than once follows the usual Wordle rule: greens are marked first, then each remaining guessed letter
// is yellow, left to right, only while the answer still has an unmatched copy of it. So guessing SPEED
// against ABIDE gives --Y-Y: the answer has one E, which goes to the first E in the guess.
//
// The reward of a guess is greens/5 + yellows/10, the same scale as Mastermind's black and white pegs,
// so a solving guess is worth 1.0. The engine adds 1 - reward to the pain total on every turn.
package wordle

import (
	"fmt"

	"github.com/chrisconley/go-hangman/internal/game"
)

const WordLength = 5

var patterns = buildPatterns()

func buildPatterns() []string {
	out := make([]string, 243)
	symbols := []byte{'-', 'Y', 'G'}
	for code := 0; code < 243; code++ {
		b := make([]byte, WordLength)
		c := code
		for i := WordLength - 1; i >= 0; i-- {
			b[i] = symbols[c%3]
			c /= 3
		}
		out[code] = string(b)
	}
	return out
}

// Score returns the response to guess when the answer is answer, as a number from 0 to 242 in which each
// letter is one base-3 digit (gray 0, yellow 1, green 2, first letter first), and the numbers of greens
// and yellows.
func Score(answer string, guess string) (code int, greens int, yellows int) {
	var remaining [26]int8
	var marks [WordLength]int8
	for i := 0; i < WordLength; i++ {
		if guess[i] == answer[i] {
			marks[i] = 2
			greens++
		} else {
			remaining[answer[i]-'a']++
		}
	}
	for i := 0; i < WordLength; i++ {
		if marks[i] == 2 {
			continue
		}
		letter := guess[i] - 'a'
		if remaining[letter] > 0 {
			marks[i] = 1
			remaining[letter]--
			yellows++
		}
	}
	for i := 0; i < WordLength; i++ {
		code = code*3 + int(marks[i])
	}
	return code, greens, yellows
}

// Feedback returns the response to guess when the answer is answer, e.g. Feedback("abbey", "babes") == "YYGG-".
func Feedback(answer string, guess string) string {
	code, _, _ := Score(answer, guess)
	return patterns[code]
}

func Reward(greens int, yellows int) float64 {
	return float64(greens)/WordLength + float64(yellows)/(2*WordLength)
}

func Act(gameplay game.GamePlay, world game.World, action game.Action) game.Reaction {
	code, greens, yellows := Score(world.String(), action.String())
	return game.NewReaction(action, patterns[code], Reward(greens, yellows))
}

// NewAgreement returns the rules of Wordle with the given words as the allowed guesses. A guess that has
// already been played in the game is not offered again.
func NewAgreement(guesses []string) game.Agreement {
	actions := make([]game.Action, len(guesses))
	for i, guess := range guesses {
		if !isWord(guess) {
			panic(fmt.Sprintf("wordle: %q is not five lowercase letters", guess))
		}
		actions[i] = game.NewAction(guess)
	}
	getActions := func(gameplay game.GamePlay, remainingWorlds game.PossibleWorlds, aa game.AvailableActions) []game.Action {
		available := make([]game.Action, 0, len(actions))
		for _, action := range actions {
			if !gameplay.PreviouslyActed(action) {
				available = append(available, action)
			}
		}
		return available
	}
	return game.Agreement{GetActions: getActions, Act: Act, NewExplore: newTableExplore(guesses)}
}

func isWord(s string) bool {
	if len(s) != WordLength {
		return false
	}
	for i := 0; i < WordLength; i++ {
		if s[i] < 'a' || s[i] > 'z' {
			return false
		}
	}
	return true
}
