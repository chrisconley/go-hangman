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
	fmt.Printf("\n%s", gs.key)
}

func (gs *GameState) Key() string {
	return gs.key
}

type DecisionMemory interface {
	Add(gs GameState, w []World)
	Get(gs GameState) ([]World, bool)
}

type NoMemory struct {
}

func (nm *NoMemory) Add(gs GameState, w []World) {

}

func (nm *NoMemory) Get(gs GameState) ([]World, bool) {
	return []World{}, false
}
