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

var comparedMotivations = []game.Motivation{
	{Entropy: 1},
	{MostParts: 1},
	{MiniMax: 1},
	{Reward: 1},
	{Entropy: 0.7, MiniMax: 0.15, MostParts: 0.15},
	{Reward: 0.75, Entropy: 0.25},
	{Reward: 0.9, MiniMax: 0.1},
	{Reward: 0.5, MostParts: 0.5},
}

// The caches of the one-at-a-time loop identify a position by the size of the word list and the guesses
// and responses so far, not by the words, so they are cleared before each list.
//
// Playing the games together must give every game exactly what it gets when the games are played one at
// a time: the same guesses, responses, turn count, reward and pain.
func TestPlayingTogetherMatchesOneAtATime(t *testing.T) {
	nineLetters := utils.ReadLinesFromFilePath("../../fixtures/nine-letter-words.txt")
	fiveLetters := utils.ReadLinesFromFilePath("../../fixtures/five-letter-words.txt")
	codes := utils.ReadLinesFromFilePath("../../fixtures/mastermind-ABCD-3.txt")
	repeat := func(lines []string, i int) []string { return append(append([]string{}, lines...), lines[i]) }
	cases := []struct {
		name      string
		agreement game.Agreement
		lines     []string
		limit     int
	}{
		{"hangman", hangman.Agreement, nineLetters[:300], 300},
		{"hangman, first 40 of 300", hangman.Agreement, nineLetters[:300], 40},
		{"hangman, a word twice", hangman.Agreement, repeat(nineLetters[:100], 7), 101},
		{"hangman, one word", hangman.Agreement, nineLetters[:1], 1},
		{"hangman, a word twice before the guess that wins", hangman.Agreement, []string{"ace", "ace", "abe"}, 3},
		{"mastermind", mastermind.Agreement, codes, len(codes)},
		{"mastermind, first 10", mastermind.Agreement, codes, 10},
		{"mastermind, a code twice", mastermind.Agreement, repeat(codes, 5), len(codes) + 1},
		{"wordle", wordle.NewAgreement(fiveLetters), fiveLetters[:120], 120},
		{"wordle, first 15 of 120", wordle.NewAgreement(fiveLetters), fiveLetters[:120], 15},
		{"wordle, a word twice", wordle.NewAgreement(fiveLetters), repeat(fiveLetters[:60], 3), 61},
		{"wordle, one word", wordle.NewAgreement(fiveLetters), fiveLetters[:1], 1},
	}
	for _, c := range cases {
		for _, motivation := range comparedMotivations {
			game.ThePredictionsFreeCache.Clear()
			game.ObserveCache.Cache.Clear()
			oneAtATime := playGamesOneAtATime(c.agreement, motivation, c.lines, c.limit)
			together := playGamesTogether(c.agreement, motivation, c.lines, c.limit)
			if len(together) != len(oneAtATime) {
				t.Fatalf("%s %+v: %d games together, %d one at a time", c.name, motivation, len(together), len(oneAtATime))
			}
			for i := range oneAtATime {
				if !reflect.DeepEqual(together[i], oneAtATime[i]) {
					t.Errorf("%s %+v: game %d (%s) differs\ntogether:       %+v\none at a time: %+v", c.name, motivation, i, c.lines[i], together[i], oneAtATime[i])
					break
				}
			}
		}
	}
}

// PlayGames guesses the last answer as soon as it is the only one left, without scoring. At every
// position with one answer left, scoring every guess the old way must pick that answer.
func TestScoringPicksTheLastWord(t *testing.T) {
	if testing.Short() {
		t.Skip("plays every game one at a time; about 12 seconds")
	}
	nineLetters := utils.ReadLinesFromFilePath("../../fixtures/nine-letter-words.txt")
	fiveLetters := utils.ReadLinesFromFilePath("../../fixtures/five-letter-words.txt")
	cases := []struct {
		name      string
		agreement game.Agreement
		lines     []string
		limit     int
	}{
		{"hangman", hangman.Agreement, nineLetters[:300], 300},
		{"hangman, a word twice", hangman.Agreement, append(append([]string{}, nineLetters[:100]...), nineLetters[7]), 101},
		{"mastermind, 3 pegs", mastermind.Agreement, utils.ReadLinesFromFilePath("../../fixtures/mastermind-ABCD-3.txt"), 64},
		{"mastermind, first 100 of 1,296", mastermind.Agreement, utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt"), 100},
		{"wordle", wordle.NewAgreement(fiveLetters), fiveLetters[:120], 120},
	}
	for _, c := range cases {
		for _, motivation := range comparedMotivations {
			checked := 0
			agreement := c.agreement
			agreement.GetActions = func(gameplay game.GamePlay, possibleWorlds game.PossibleWorlds, aa game.AvailableActions) []game.Action {
				actions := c.agreement.GetActions(gameplay, possibleWorlds, aa)
				if possibleWorlds.Len() == 1 {
					checked++
					word := game.BUniverse.GetWorldByIndex(possibleWorlds.Indices()[0])
					picked := game.Decide(game.ExploreReactions(gameplay, possibleWorlds, actions, c.agreement.Act), motivation)
					if !picked.Action.Matches(word) {
						t.Errorf("%s %+v: with only %s left, scoring picks %s", c.name, motivation, word, picked.Action)
					}
				}
				return actions
			}
			game.ThePredictionsFreeCache.Clear()
			game.ObserveCache.Cache.Clear()
			playGamesOneAtATime(agreement, motivation, c.lines, c.limit)
			if checked == 0 {
				t.Errorf("%s %+v: no position with one answer left", c.name, motivation)
			}
		}
	}
}
