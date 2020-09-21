package game

type AvailableActions struct {
	worlds       []World
	actions      []Action
	worldActions []Action
	worldSize    int
}

func BuildAvailableActions(worlds PossibleWorlds) AvailableActions {
	aa := AvailableActions{
		actions:      []Action{},
		worldActions: []Action{},
		worldSize:    len(worlds.Worlds[0]),
	}
	keys := make(map[Action]bool)
	for _, world := range worlds.Worlds {
		worldAction := Action(world)
		aa.worldActions = append(aa.worldActions, worldAction)
		for i := 0; i < len(worldAction); i++ {
			action := Action(worldAction[i])
			if _, ok := keys[action]; !ok {
				keys[action] = true
				aa.actions = append(aa.actions, action)
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
