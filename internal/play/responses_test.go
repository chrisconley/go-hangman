package play

import (
	"reflect"
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/chrisconley/go-hangman/internal/wordle"
)

// At every position the walk reaches, the scores from the table of responses must be exactly the scores
// game.ExploreReactions works out from the rules, guess for guess and in the same order.
func TestTableScoresMatchTheEngine(t *testing.T) {
	nineLetters := utils.ReadLinesFromFilePath("../../fixtures/nine-letter-words.txt")
	fiveLetters := utils.ReadLinesFromFilePath("../../fixtures/five-letter-words.txt")
	commercial := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")
	cases := []struct {
		name      string
		agreement game.Agreement
		lines     []string
	}{
		{"hangman", hangman.Agreement, nineLetters[:400]},
		{"hangman, a word twice", hangman.Agreement, append(append([]string{}, nineLetters[:100]...), nineLetters[7])},
		{"mastermind", mastermind.Agreement, commercial[:400]},
		{"mastermind, 3 pegs", mastermind.Agreement, utils.ReadLinesFromFilePath("../../fixtures/mastermind-ABCD-3.txt")},
		{"wordle", wordle.NewAgreement(fiveLetters), fiveLetters[:150]},
	}
	for _, c := range cases {
		for _, motivation := range []game.Motivation{{Entropy: 1}, {MostParts: 1}, {Reward: 0.75, Entropy: 0.25}} {
			positions := 0
			check := func(gameplay game.GamePlay, possibleWorlds game.PossibleWorlds, actions []game.Action, predictions []game.Prediction) {
				positions++
				want := game.ExploreReactions(gameplay, possibleWorlds, actions, c.agreement.Act).PredictionsMap.Values
				if !reflect.DeepEqual(predictions, want) {
					t.Fatalf("%s %+v: scores differ after %d guesses\ntable: %+v\nrules: %+v", c.name, motivation, gameplay.TurnCount, predictions, want)
				}
			}
			game.BUniverse = game.InitUniverse(c.lines)
			w := walk{agreement: c.agreement, motivation: motivation, limit: uint64(len(c.lines)), games: make([]game.GamePlay, len(c.lines)), table: newResponseTable(), checkScores: check}
			all := responses{answers: make([]uint64, len(c.lines))}
			for i := range all.answers {
				all.answers[i] = uint64(i)
			}
			w.play(game.InitGamePlay(game.BUniverse.PossibleWorlds.GetHash()), game.BUniverse.PossibleWorlds, all)
			if positions == 0 {
				t.Fatalf("%s: no positions checked", c.name)
			}
		}
	}
}
