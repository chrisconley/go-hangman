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
// Guesses are scored from a table of responses (see responseTable) instead of with game.ExploreReactions.
// The scores are the same numbers.
//
// As when playing one game at a time, every line is a possible answer, but only the games for the first
// limit lines are played: a group is followed only while it holds one of those answers.
func playGamesTogether(agreement game.Agreement, motivation game.Motivation, lines []string, limit int) []game.GamePlay {
	game.BUniverse = game.InitUniverse(lines)
	w := walk{agreement: agreement, motivation: motivation, limit: uint64(limit), games: make([]game.GamePlay, len(lines)), table: newResponseTable()}
	all := responses{answers: make([]uint64, len(lines))}
	for i := range all.answers {
		all.answers[i] = uint64(i)
	}
	root := game.InitGamePlay(game.BUniverse.PossibleWorlds.GetHash())
	w.play(root, game.BUniverse.PossibleWorlds, all)
	return w.games[:limit]
}

type walk struct {
	agreement  game.Agreement
	motivation game.Motivation
	limit      uint64
	games      []game.GamePlay
	finished   []game.GamePlay
	table      *responseTable
	// checkScores, when set by a test, is shown the scores at every position.
	checkScores func(gameplay game.GamePlay, possibleWorlds game.PossibleWorlds, actions []game.Action, predictions []game.Prediction)
}

// play scores the guesses at a position, picks one, and carries on with each group of answers that give
// the same response to it. r holds the answers still possible, in list order, and the rows made so far.
func (w *walk) play(gameplay game.GamePlay, possibleWorlds game.PossibleWorlds, r responses) {
	if len(r.answers) == 1 {
		w.guessTheLastWord(gameplay, r.answers[0])
		return
	}
	actions := w.agreement.GetActions(gameplay, possibleWorlds, game.BUniverse.Actions)
	r, rowsOfActions, rowCount := w.table.add(gameplay, r, actions, w.agreement.Act)
	defer w.table.drop(rowCount)
	predictions := w.table.score(r, rowsOfActions, possibleWorlds, actions)
	if w.checkScores != nil {
		w.checkScores(gameplay, possibleWorlds, actions, predictions)
	}
	logPredictions(predictions)
	action := game.DecideAction(predictions, w.motivation)
	var row int32
	for i, prediction := range predictions {
		if prediction.Action.Matches(action) {
			logDecision(prediction)
		}
		if actions[i].Matches(action) {
			row = rowsOfActions[i]
		}
	}

	numbers := r.row(row)
	groups := make([][]int, len(w.table.rows[row].rewards))
	order := []uint16{}
	for place, number := range numbers {
		if groups[number] == nil {
			order = append(order, number)
		}
		groups[number] = append(groups[number], place)
	}

	for _, number := range order {
		places := groups[number]
		first := r.answers[places[0]]
		if first >= w.limit {
			continue
		}
		reaction := w.agreement.Act(gameplay, game.BUniverse.GetWorldByIndex(first), action)
		next := gameplay
		next.Reactions = append(make([]game.Reaction, 0, len(gameplay.Reactions)+1), gameplay.Reactions...)
		next.AddTurn(reaction)
		logReaction(reaction)
		group := w.table.keep(r, places)
		if action.Matches(game.BUniverse.GetWorldByIndex(first)) {
			w.finish(group.answers, action, next)
			continue
		}
		nextWorlds := game.NewPossibleWorlds(game.BUniverse.Size)
		for _, answer := range group.answers {
			nextWorlds.Add(answer)
		}
		logPossibleWorlds(nextWorlds)
		w.play(next, nextWorlds, group)
	}
}

// guessTheLastWord plays the one answer left, without scoring any guesses. Scoring would pick it anyway:
// with one answer left every guess leaves one group, so every guess has entropy 0, MiniMax 0 and one
// part, and the answer itself has the top reward, 1. Whatever the weights, the answer ties for the top
// weight, and a tie goes to a guess that could be the answer, which only the answer itself is. The
// answer is played even when the game's GetActions would not have offered it.
func (w *walk) guessTheLastWord(gameplay game.GamePlay, answer uint64) {
	truth := game.BUniverse.GetWorldByIndex(answer)
	action := game.NewAction(truth.String())
	reaction := w.agreement.Act(gameplay, truth, action)
	next := gameplay
	next.Reactions = append(make([]game.Reaction, 0, len(gameplay.Reactions)+1), gameplay.Reactions...)
	next.AddTurn(reaction)
	logReaction(reaction)
	w.finish([]uint64{answer}, action, next)
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
