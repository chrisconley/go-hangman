package game

import "sort"

type PossibleWorlds struct {
	Worlds []World
	key    string
}

func BuildPossibleWorlds(worlds []World) PossibleWorlds {
	sort.SliceStable(worlds, func(i, j int) bool {
		return worlds[i] < worlds[j]
	})
	key := ""
	for _, world := range worlds {
		key += string(world)
	}
	return PossibleWorlds{Worlds: worlds, key: key}
}

func (w *PossibleWorlds) Key() string {
	return w.key
}

func (pw *PossibleWorlds) Len() int {
	return len(pw.Worlds)
}
