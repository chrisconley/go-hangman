package game

import "fmt"

// type geometry interface {
//     area() float64
//     perim() float64
// }

type GameState struct {
	reactions []Reaction
	key       string
}

func (gs *GameState) Add(reaction Reaction) {
	gs.reactions = append(gs.reactions, reaction)
	gs.key = fmt.Sprintf("%s|%v:%v", gs.key, reaction.Action, reaction.Response)
}

func (gs *GameState) Key() string {
	return gs.key
}

type DecisionMemory interface {
	Init()
	Add(gs GameState, w PossibleWorlds)
	Get(gs GameState) (PossibleWorlds, bool)
}

func InitMemory(m DecisionMemory) DecisionMemory {
	m.Init()
	return m
}

type NoMemory struct {
}

func (nm *NoMemory) Init() {
}

func (nm *NoMemory) Add(gs GameState, w PossibleWorlds) {
}

func (nm *NoMemory) Get(gs GameState) (PossibleWorlds, bool) {
	return BuildPossibleWorlds([]string{}), false
}

type NaiveMemory struct {
	memory map[string]PossibleWorlds
}

func (nm *NaiveMemory) Init() {
	nm.memory = make(map[string]PossibleWorlds)
}

func (nm *NaiveMemory) Add(gs GameState, w PossibleWorlds) {
	nm.memory[gs.Key()] = w

}

func (nm *NaiveMemory) Get(gs GameState) (PossibleWorlds, bool) {
	if worlds, ok := nm.memory[gs.Key()]; ok {
		return worlds, true
	} else {
		return BuildPossibleWorlds([]string{}), false
	}
}
