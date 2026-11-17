package wordle

import (
	"reflect"
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/play"
	"github.com/chrisconley/go-hangman/internal/utils"
)

// Scoring from the response table must give exactly the numbers game.ExploreReactions gives, and the
// same chosen guess, at every position the games reach.
func TestTableScoresMatchTheEngine(t *testing.T) {
	words := utils.ReadLinesFromFilePath("../../fixtures/five-letter-words.txt")
	answers := words[:150]
	agreement := NewAgreement(words)
	checked := 0
	for _, motivation := range []game.Motivation{{Entropy: 1}, {MostParts: 1}, {MiniMax: 1}, {Reward: 0.75, Entropy: 0.25}} {
		games := play.PlayGames(agreement, motivation, answers, len(answers))
		universe := game.InitUniverse(answers)
		game.BUniverse = universe
		table := agreement.NewExplore(universe)
		seen := map[uint64]bool{}
		for _, gameplay := range games {
			position := game.InitGamePlay(universe.PossibleWorlds.GetHash())
			worlds := universe.PossibleWorlds
			for _, reaction := range gameplay.Reactions {
				if !seen[position.Hash] {
					seen[position.Hash] = true
					checked++
					actions := agreement.GetActions(position, worlds, universe.Actions)
					fromTable := table(position, worlds, actions)
					fromEngine := game.ExploreReactions(position, worlds, actions, Act)
					for _, action := range actions {
						if a, b := fromTable.Get(action.String()), fromEngine.Get(action.String()); !reflect.DeepEqual(a, b) {
							t.Fatalf("%+v, after %d guesses, %s: table %+v, engine %+v", motivation, position.TurnCount, action.String(), a, b)
						}
					}
					if a, b := game.Decide(fromTable, motivation).Action.String(), game.Decide(fromEngine, motivation).Action.String(); a != b {
						t.Fatalf("%+v, after %d guesses: table chooses %s, engine %s", motivation, position.TurnCount, a, b)
					}
				}
				next := game.NewPossibleWorlds(universe.Size)
				for _, w := range worlds.Indices() {
					if Act(position, universe.GetWorldByIndex(w), reaction.Action).Response.Matches(reaction.Response) {
						next.Add(w)
					}
				}
				worlds = next
				position.AddTurn(reaction)
			}
		}
	}
	if checked < 100 {
		t.Errorf("only %d positions checked", checked)
	}
}
