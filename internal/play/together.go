package play

import (
	"github.com/chrisconley/go-hangman/internal/game"
)

// playGamesTogether plays the same games as playing them one at a time, but all at once. A position
// (the guesses so far and their responses) is scored once, however many games reach it: the answers
// still possible are split by the response they give to the chosen guess, and each group carries on
// from there. The guess chosen at a position depends only on the position, so every game gets the same
// guesses, responses and rewards as when it is played on its own.
//
// As when playing one game at a time, every line is a possible answer, but only the games for the first
// limit lines are played: a group is followed only while it holds one of those answers.
func playGamesTogether(agreement game.Agreement, motivation game.Motivation, lines []string, limit int) []game.GamePlay {
	game.BUniverse = game.InitUniverse(lines)
	w := walk{agreement: agreement, motivation: motivation, limit: uint64(limit), games: make([]game.GamePlay, len(lines))}
	root := game.InitGamePlay(game.BUniverse.PossibleWorlds.GetHash())
	w.play(root, game.BUniverse.PossibleWorlds)
	return w.games[:limit]
}

type walk struct {
	agreement  game.Agreement
	motivation game.Motivation
	limit      uint64
	games      []game.GamePlay
	finished   []game.GamePlay
}

func (w *walk) play(gameplay game.GamePlay, possibleWorlds game.PossibleWorlds) {
	actions := w.agreement.GetActions(gameplay, possibleWorlds, game.BUniverse.Actions)
	predictions := game.ExploreReactions(gameplay, possibleWorlds, actions, w.agreement.Act)
	logPredictions(predictions)
	action := Decide(predictions, w.motivation)

	groups := map[string][]uint64{}
	reactions := map[string]game.Reaction{}
	order := []string{}
	for _, worldIndex := range possibleWorlds.Indices() {
		reaction := w.agreement.Act(gameplay, game.BUniverse.GetWorldByIndex(worldIndex), action)
		response := reaction.Response.String()
		if first, ok := reactions[response]; !ok {
			reactions[response] = reaction
			order = append(order, response)
		} else if first.Reward != reaction.Reward {
			panic("two answers gave the same response but different rewards; playing games together needs the reward to follow from the response")
		}
		groups[response] = append(groups[response], worldIndex)
	}

	for _, response := range order {
		group := groups[response]
		if group[0] >= w.limit {
			continue
		}
		next := gameplay
		next.Reactions = append(make([]game.Reaction, 0, len(gameplay.Reactions)+1), gameplay.Reactions...)
		next.AddTurn(reactions[response])
		logReaction(reactions[response])
		if action.Matches(game.BUniverse.GetWorldByIndex(group[0])) {
			w.finish(group, action, next)
			continue
		}
		nextWorlds := game.NewPossibleWorlds(game.BUniverse.Size)
		for _, worldIndex := range group {
			nextWorlds.Add(worldIndex)
		}
		logPossibleWorlds(nextWorlds)
		w.play(next, nextWorlds)
	}
}

// finish ends the game of every answer in a group the guess has just solved. A group holds more than
// one answer only when the same word appears more than once in the list.
func (w *walk) finish(group []uint64, action game.Action, gameplay game.GamePlay) {
	for _, worldIndex := range group {
		truth := game.BUniverse.GetWorldByIndex(worldIndex)
		if !action.Matches(truth) {
			panic("a guess solved one answer in a group but not another")
		}
		if worldIndex >= w.limit {
			continue
		}
		w.games[worldIndex] = gameplay
		w.finished = append(w.finished, gameplay)
		logProgress(uint64(len(w.finished)), truth, w.finished)
	}
}
