package game

type AvailableActions struct {
	worlds  []World
	actions []Action
}

func BuildAvailableActions(worlds []World) AvailableActions {
	aa := AvailableActions{actions: []Action{}}
	keys := make(map[Action]bool)
	for _, world := range worlds {
		worldAction := Action(world)
		aa.actions = append(aa.actions, worldAction)
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

func (aa *AvailableActions) All() []Action {
	return aa.actions
}

func (aa *AvailableActions) Update(usedAction Action) []Action {
	actions := []Action{}
	for _, action := range aa.actions {
		if action != usedAction {
			actions = append(actions, action)
		}
	}
	aa.actions = actions
	return aa.actions
}
