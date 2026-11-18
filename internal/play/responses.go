package play

import (
	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/mth"
)

// responseTable scores guesses from a table of responses. For each guess it holds the response every
// answer still possible gives to it, as a small number, so a guess is scored by counting numbers instead
// of working the responses out from the rules again at every position.
//
// A guess gets a row the first time it is offered on the way down from the first position. The row is
// worked out with the game's Act for the answers still possible there, and the responses are numbered in
// the order they first appear. Each position below keeps a copy of every row, cut down to its own
// answers, so the numbers for one position sit together in memory. When the walk goes back up past the
// position that made a row, the row goes away. This assumes, as all three games do, that the response
// to a guess does not depend on the guesses played before it.
//
// So rows exist only for guesses offered on the current path: every allowed guess for Wordle, the
// letters and a few whole words for hangman (words are offered only when few answers remain), and the
// guesses Mastermind's GetActions offers.
type responseTable struct {
	guessNumbers map[uint64]int32
	// rowOf[g] is the row of guess number g on the current path, or -1.
	rowOf []int32
	rows  []responseRow
	// scratch space, reused at every position
	responseNumbers map[uint64]uint16
	counts          []int32
	firstSeen       []uint16
	groupSizes      []float64
	groupRewards    []float64
	groupRewardList []float64
	position        int
}

type responseRow struct {
	guess int32
	// rewards[n] is the reward that goes with response number n.
	rewards []float64
	// solvedBy lists the answers that are the guess itself.
	solvedBy []uint64
	// scoredAt stops a guess offered twice at one position from being scored twice.
	scoredAt int
}

// responses holds the response numbers of every row for the answers still possible at one position:
// row r, answer j is numbers[r*len(answers)+j].
type responses struct {
	answers []uint64
	numbers []uint16
}

func (r responses) row(row int32) []uint16 {
	n := len(r.answers)
	return r.numbers[int(row)*n : (int(row)+1)*n]
}

func newResponseTable() *responseTable {
	return &responseTable{
		guessNumbers:    map[uint64]int32{},
		responseNumbers: map[uint64]uint16{},
		counts:          make([]int32, 1<<16),
		groupRewards:    make([]float64, 1<<16),
	}
}

// add gives a row to every guess offered at this position that has none yet. It returns the responses
// with the new rows added, the row of each guess offered, and the number of rows there were before, for
// drop.
func (t *responseTable) add(gameplay game.GamePlay, r responses, actions []game.Action, act game.Act) (responses, []int32, int) {
	before := len(t.rows)
	rowsOfActions := make([]int32, len(actions))
	added := []game.Action{}
	for i, action := range actions {
		guess, ok := t.guessNumbers[action.Hash()]
		if !ok {
			guess = int32(len(t.rowOf))
			t.guessNumbers[action.Hash()] = guess
			t.rowOf = append(t.rowOf, -1)
		}
		if t.rowOf[guess] < 0 {
			t.rowOf[guess] = int32(len(t.rows))
			t.rows = append(t.rows, responseRow{guess: guess})
			added = append(added, action)
		}
		rowsOfActions[i] = t.rowOf[guess]
	}
	if len(added) == 0 {
		return r, rowsOfActions, before
	}

	n := len(r.answers)
	numbers := make([]uint16, len(t.rows)*n)
	copy(numbers, r.numbers)
	r = responses{answers: r.answers, numbers: numbers}
	for i, action := range added {
		row := int32(before + i)
		numbers := r.row(row)
		info := &t.rows[row]
		for response := range t.responseNumbers {
			delete(t.responseNumbers, response)
		}
		for j, answer := range r.answers {
			world := game.BUniverse.GetWorldByIndex(answer)
			reaction := act(gameplay, world, action)
			number, ok := t.responseNumbers[reaction.Response.Hash()]
			if !ok {
				if len(info.rewards) == len(t.counts) {
					panic("a guess has more than 65,536 different responses")
				}
				number = uint16(len(info.rewards))
				t.responseNumbers[reaction.Response.Hash()] = number
				info.rewards = append(info.rewards, reaction.Reward)
			} else if info.rewards[number] != reaction.Reward {
				panic("two answers gave the same response but different rewards; playing games together needs the reward to follow from the response")
			}
			numbers[j] = number
			if action.Matches(world) {
				info.solvedBy = append(info.solvedBy, answer)
			}
		}
	}
	return r, rowsOfActions, before
}

// drop removes the rows made after the first rowCount, when the walk leaves the position that made them.
func (t *responseTable) drop(rowCount int) {
	for _, row := range t.rows[rowCount:] {
		t.rowOf[row.guess] = -1
	}
	t.rows = t.rows[:rowCount]
}

// keep cuts every row down to the answers at the given places in r, for a position below.
func (t *responseTable) keep(r responses, places []int) responses {
	n, m := len(places), len(r.answers)
	kept := responses{answers: make([]uint64, n), numbers: make([]uint16, len(t.rows)*n)}
	for j, place := range places {
		kept.answers[j] = r.answers[place]
	}
	for row := range t.rows {
		from := r.numbers[row*m : (row+1)*m]
		to := kept.numbers[row*n : (row+1)*n]
		for j, place := range places {
			to[j] = from[place]
		}
	}
	return kept
}

// score gives every guess the same scores as game.ExploreReactions, in the same order. It counts the
// groups answer by answer and adds up each group's rewards, then takes the groups in the fixed order of
// game.GroupsInFixedOrder, as game.ExploreReactions does, so that every sum comes out the same. A guess
// offered twice is scored once, as it is in game.Predictions.
func (t *responseTable) score(r responses, rowsOfActions []int32, possibleWorlds game.PossibleWorlds, actions []game.Action) []game.Prediction {
	t.position++
	total := float64(len(r.answers))
	predictions := make([]game.Prediction, 0, len(actions))
	for i, action := range actions {
		info := &t.rows[rowsOfActions[i]]
		if info.scoredAt == t.position {
			continue
		}
		info.scoredAt = t.position

		firstSeen := t.firstSeen[:0]
		for _, number := range r.row(rowsOfActions[i]) {
			if t.counts[number] == 0 {
				firstSeen = append(firstSeen, number)
			}
			t.counts[number]++
			t.groupRewards[number] += info.rewards[number]
		}
		groupSizes := t.groupSizes[:0]
		groupRewardList := t.groupRewardList[:0]
		for _, number := range firstSeen {
			groupSizes = append(groupSizes, float64(t.counts[number]))
			groupRewardList = append(groupRewardList, t.groupRewards[number])
			t.counts[number] = 0
			t.groupRewards[number] = 0
		}
		t.firstSeen, t.groupSizes, t.groupRewardList = firstSeen, groupSizes, groupRewardList
		reward, orderedSizes := game.GroupsInFixedOrder(groupSizes, groupRewardList)

		prediction := game.Prediction{Action: action}
		prediction.Reward = reward / total
		prediction.Entropy = mth.CalcEntropy(total, orderedSizes)
		prediction.PartsCount = float64(len(groupSizes))
		prediction.InversedMax = mth.FindInversedMax(total, groupSizes)
		for _, answer := range info.solvedBy {
			if possibleWorlds.Bits.Contains(answer) {
				prediction.ActionMatchesAPossibleWorld = true
				break
			}
		}
		predictions = append(predictions, prediction)
	}
	return predictions
}
