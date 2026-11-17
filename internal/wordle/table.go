package wordle

import (
	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/mth"
)

const allGreen = 242

// newTableExplore works out, once, the response every allowed guess gets from every possible answer,
// and returns a function that scores guesses by looking responses up in that table instead of calling
// Act for each guess and answer. A guess gets exactly the numbers game.ExploreReactions gives it: the
// rewards are added in the same order (answer by answer) and the response counts are kept in the order
// the responses first appear, which is the order the entropy is summed in.
func newTableExplore(guesses []string) func(universe game.Universe) game.Explore {
	return func(universe game.Universe) game.Explore {
		row := make(map[string][]uint8, len(guesses))
		for _, guess := range guesses {
			responses := make([]uint8, len(universe.Worlds))
			for w, answer := range universe.Worlds {
				code, _, _ := Score(answer.String(), guess)
				responses[w] = uint8(code)
			}
			row[guess] = responses
		}
		var rewardOf [243]float64
		for code := 0; code < 243; code++ {
			greens, yellows := 0, 0
			for c := code; c > 0; c /= 3 {
				if c%3 == 2 {
					greens++
				} else if c%3 == 1 {
					yellows++
				}
			}
			rewardOf[code] = Reward(greens, yellows)
		}
		return func(gameplay game.GamePlay, possibleWorlds game.PossibleWorlds, actions []game.Action) game.Predictions {
			worlds := possibleWorlds.Indices()
			total := float64(len(worlds))
			var counts [243]int
			var order [243]uint8
			var groupSizes [243]float64
			predictions := game.NewPredictions()
			for _, action := range actions {
				responses := row[action.String()]
				parts := 0
				rewards := 0.0
				for _, w := range worlds {
					code := responses[w]
					if counts[code] == 0 {
						order[parts] = code
						parts++
					}
					counts[code]++
					rewards += rewardOf[code]
				}
				for i := 0; i < parts; i++ {
					groupSizes[i] = float64(counts[order[i]])
				}
				predictions.AddPrediction(action.String(), game.Prediction{
					Action:                      action,
					Reward:                      rewards / total,
					Entropy:                     mth.CalcEntropy(total, groupSizes[:parts]),
					InversedMax:                 mth.FindInversedMax(total, groupSizes[:parts]),
					PartsCount:                  float64(parts),
					ActionMatchesAPossibleWorld: counts[allGreen] > 0,
				})
				for i := 0; i < parts; i++ {
					counts[order[i]] = 0
				}
			}
			return predictions
		}
	}
}
