package game

// ###########################
// ### THE SYSTEM ###
// ###########################
type World string
type Action string
type Response string
type Agreement func(w World, action Action) (reaction Reaction)
type ActionFilter func(w Worlds, aa AvailableActions) []Action

type Reaction struct {
	Action   Action
	Response Response
	Reward   float64
}
