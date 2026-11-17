package main

import (
	"container/list"
	"reflect"
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/hangman"
	"github.com/chrisconley/go-hangman/internal/mastermind"
	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/chrisconley/go-hangman/internal/wordle"
)

var testMotivations = []game.Motivation{
	{Entropy: 1},
	{MostParts: 1},
	{MiniMax: 1},
	{Reward: 1},
	{Entropy: 0.7, MiniMax: 0.15, MostParts: 0.15},
	{Reward: 0.75, Entropy: 0.25},
}

func newSettings(mode string, lines []string) trialSettings {
	return trialSettings{mode: mode, limit: len(lines), cacheSize: 1000, cache: &stateCache{capacity: 16, order: list.New(), entries: map[uint64]*list.Element{}}, cacheDepth: 1}
}

func totals(t *testing.T, agreement game.Agreement, lines []string, motivation game.Motivation, s trialSettings) play.Stats {
	games, _ := playTrial(agreement, lines, motivation, s)
	for i, gameplay := range games {
		if gameplay.TurnCount == 0 {
			t.Fatalf("%+v %s: game %d (%s) was not played", motivation, s.mode, i, lines[i])
		}
	}
	return play.CalculateStats(games, true)
}

func assertSameStats(t *testing.T, name string, motivation game.Motivation, games play.Stats, tree play.Stats) {
	if games.TotalTurns != tree.TotalTurns || games.MaxTurns != tree.MaxTurns || games.TotalReward != tree.TotalReward || games.TotalPain != tree.TotalPain || !reflect.DeepEqual(games.TurnCounts, tree.TurnCounts) {
		t.Errorf("%s %+v: games mode %d turns (max %d, reward %v, pain %v), tree mode %d turns (max %d, reward %v, pain %v)",
			name, motivation, games.TotalTurns, games.MaxTurns, games.TotalReward, games.TotalPain, tree.TotalTurns, tree.MaxTurns, tree.TotalReward, tree.TotalPain)
	}
}

// Tree mode scores each position once and games mode plays one game at a time; both must make the same
// guesses, so every total and every game's turn count must agree.
func TestTreeModeMatchesGamesMode(t *testing.T) {
	words := utils.ReadLinesFromFilePath("../../fixtures/five-letter-words.txt")
	hangmanWords := 1000
	if testing.Short() {
		hangmanWords = 200
	}
	cases := []struct {
		name      string
		agreement game.Agreement
		lines     []string
	}{
		{"hangman", hangman.Agreement, utils.ReadLinesFromFilePath("../../fixtures/nine-letter-words.txt")[:hangmanWords]},
		{"mastermind", mastermind.Agreement, utils.ReadLinesFromFilePath("../../fixtures/mastermind-ABCD-3.txt")},
		{"mastermind 4x6", mastermind.Agreement, utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")},
		{"wordle", wordle.NewAgreement(words), words[:150]},
	}
	for _, c := range cases {
		if testing.Short() && c.name == "mastermind 4x6" {
			continue
		}
		for _, motivation := range testMotivations {
			gamesStats := totals(t, c.agreement, c.lines, motivation, newSettings("games", c.lines))
			treeStats := totals(t, c.agreement, c.lines, motivation, newSettings("tree", c.lines))
			assertSameStats(t, c.name, motivation, gamesStats, treeStats)
		}
	}
}

// Scoring Wordle guesses from the table of responses (-tree.fast) must give the same scores and pick the
// same guess as the engine at every position; -tree.check stops with an error on any difference.
func TestWordleTableMatchesEngine(t *testing.T) {
	words := utils.ReadLinesFromFilePath("../../fixtures/five-letter-words.txt")
	answers := words[:200]
	agreement := wordle.NewAgreement(words)
	for _, motivation := range testMotivations {
		slow := totals(t, agreement, answers, motivation, newSettings("tree", answers))
		s := newSettings("tree", answers)
		s.fast = newFastWordle(words, game.InitUniverse(answers).Worlds)
		s.check = true
		games, checked := playTrial(agreement, answers, motivation, s)
		if checked == 0 {
			t.Errorf("%+v: no positions were checked", motivation)
		}
		assertSameStats(t, "wordle table", motivation, slow, play.CalculateStats(games, true))
	}
}

// Standard Mastermind (4 pegs, 6 colours) from the README's Results table.
func TestMastermindResults(t *testing.T) {
	if testing.Short() {
		t.Skip("about 10 seconds")
	}
	lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")
	for _, c := range []struct {
		motivation game.Motivation
		turns      int
		maxTurns   int
	}{
		{game.Motivation{MostParts: 1}, 5668, 6},
		{game.Motivation{Entropy: 1}, 5723, 6},
	} {
		stats := totals(t, mastermind.Agreement, lines, c.motivation, newSettings("tree", lines))
		if stats.TotalTurns != c.turns || stats.MaxTurns != c.maxTurns {
			t.Errorf("%+v: %d turns, at most %d; want %d, at most %d", c.motivation, stats.TotalTurns, stats.MaxTurns, c.turns, c.maxTurns)
		}
	}
}

// Nine-letter hangman, Entropy alone, from the README's Results table.
func TestHangmanEntropyResult(t *testing.T) {
	lines := utils.ReadLinesFromFilePath("../../fixtures/nine-letter-words.txt")
	stats := totals(t, hangman.Agreement, lines, game.Motivation{Entropy: 1}, newSettings("tree", lines))
	if stats.TotalTurns != 146356 || stats.MaxTurns != 13 {
		t.Errorf("%d turns, at most %d; want 146,356, at most 13", stats.TotalTurns, stats.MaxTurns)
	}
}

func TestBuildTrialsAndShard(t *testing.T) {
	trials := buildTrials(0.5, "0.75:0.25:0:0", false)
	if len(trials) != 11 || trials[10] != (game.Motivation{Reward: 0.75, Entropy: 0.25}) {
		t.Fatalf("expected the 10 combinations of step 0.5 and the listed one, got %v", trials)
	}
	if n := len(buildTrials(0.5, "", true)); n != 6 {
		t.Errorf("-mostparts.zero at step 0.5: %d combinations, want 6", n)
	}
	seen := map[game.Motivation]int{}
	for _, contiguous := range []bool{false, true} {
		for index := 0; index < 3; index++ {
			for _, m := range shard(trials, index, 3, contiguous) {
				seen[m]++
			}
		}
	}
	for _, m := range trials {
		if seen[m] != 2 {
			t.Errorf("%+v is given to %d processes over the two ways of splitting, want 2", m, seen[m])
		}
	}
}
