package game

import "github.com/segmentio/fasthash/fnv1a"

type Act func(gp GamePlay, w World, action Action) (reaction Reaction)
type GetActions func(gp GamePlay, w PossibleWorlds, aa AvailableActions) []Action

type Agreement struct {
	GetActions GetActions
	Act        Act
}

type Reaction struct {
	Action   Action
	Response Response
	Reward   float64
}

type GamePlay struct {
	Reactions   []Reaction
	TurnCount   int
	TotalReward float64
	TotalPain   float64
	Hash        uint64
}

func InitGamePlay() GamePlay {
	gp := GamePlay{}
	gp.Hash = fnv1a.Init64
	return gp
}

func (gp *GamePlay) AddTurn(reaction Reaction) {
	gp.Reactions = append(gp.Reactions, reaction)
	gp.TurnCount++
	gp.TotalReward += reaction.Reward
	gp.TotalPain += 1.0 - reaction.Reward
	gp.Hash = fnv1a.AddBytes64(gp.Hash, reaction.Action.Bytes())
}

func (gp *GamePlay) PreviouslyActed(action Action) bool {
	for _, reaction := range gp.Reactions {
		if reaction.Action.Matches(action) {
			return true
		}
	}
	return false
}
