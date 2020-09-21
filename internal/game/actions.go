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
		worldActions: make([]Action, worlds.Len()),
		worldSize:    len(worlds.Worlds[0]),
	}
	// keys := make(map[Action]bool)
	for index, world := range worlds.Worlds {
		worldAction := Action(world)
		aa.worldActions[index] = worldAction
		// for i := 0; i < len(worldAction); i++ {
		// 	action := Action(worldAction[i])
		// 	if _, ok := keys[action]; !ok {
		// 		keys[action] = true
		// 		aa.actions = append(aa.actions, action)
		// 	}

		// }
	}

	aa.actions = []Action{
		Action("a"),
		Action("b"),
		Action("c"),
		Action("d"),
		Action("e"),
		Action("f"),
		Action("g"),
		Action("h"),
		Action("i"),
		Action("j"),
		Action("k"),
		Action("l"),
		Action("m"),
		Action("n"),
		Action("o"),
		Action("p"),
		Action("q"),
		Action("r"),
		Action("s"),
		Action("t"),
		Action("u"),
		Action("v"),
		Action("w"),
		Action("x"),
		Action("y"),
		Action("z"),
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
