package mastermind_test

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/utils"
)

// Before the first guess of 3 pegs and 4 colours, every code with three different colours splits the 64
// answers into groups of exactly the same sizes and rewards, so each must get exactly the same scores.
// Adding the rewards answer by answer gave DAB a reward one rounding step above ABC's.
func TestGuessesWithTheSameGroupsScoreExactlyTheSame(t *testing.T) {
	lines := utils.ReadLinesFromFilePath("../../fixtures/mastermind-ABCD-3.txt")
	game.BUniverse = game.InitUniverse(lines)
	gameplay := game.InitGamePlay(game.BUniverse.PossibleWorlds.GetHash())
	guesses := []game.Action{}
	for _, line := range lines {
		if line[0] != line[1] && line[0] != line[2] && line[1] != line[2] {
			guesses = append(guesses, game.NewAction(line))
		}
	}
	scores := game.ExploreReactions(gameplay, game.BUniverse.PossibleWorlds, guesses, mastermind.Act)
	first := scores.Get(guesses[0].String())
	for _, guess := range guesses[1:] {
		score := scores.Get(guess.String())
		if score.Reward != first.Reward || score.Entropy != first.Entropy || score.InversedMax != first.InversedMax || score.PartsCount != first.PartsCount {
			t.Errorf("%s scores %+v, %s scores %+v", guesses[0], first, guess, score)
		}
	}
}
