package game

type AvailableActions struct {
	worlds       []World
	actions      []Action
	worldActions []Action
	worldSize    int
}

// TODO: Just move this to universe creation then the performance shouldn't matter
// because it'll get called once instead of every new game with NewPerspective
func BuildAvailableActions(worlds PossibleWorlds) AvailableActions {
	aa := AvailableActions{
		actions:      []Action{},
		worldActions: make([]Action, worlds.Len()),
		worldSize:    len(worlds.Worlds[0]),
	}

	keys := make(map[int32]bool)
	for index, world := range worlds.Worlds {
		worldAction := Action(world)
		aa.worldActions[index] = worldAction
		for _, runeValue := range string(world) {
			if _, ok := keys[runeValue]; !ok {
				keys[runeValue] = true
				aa.actions = append(aa.actions, Action(runeValue))
			}
		}
	}

	return aa
}

func (aa *AvailableActions) AllActions() []Action {
	return append(aa.actions, aa.worldActions...)
}

func (aa *AvailableActions) Actions() []Action {
	return append(aa.actions)
}

func (aa *AvailableActions) WorldActions() []Action {
	return append(aa.worldActions)
}

func (aa *AvailableActions) WorldActionsLength() int {
	return len(aa.worldActions)
}

func (aa *AvailableActions) ActionsLength() int {
	return len(aa.actions)
}
