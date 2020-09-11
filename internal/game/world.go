package game

import "sort"

type PossibleWorlds struct {
	Worlds []World
	key    string
}

func BuildPossibleWorlds(lines []string) PossibleWorlds {
	sort.SliceStable(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})
	key := ""
	worlds := []World{}
	for _, line := range lines {
		key += string(line)
		worlds = append(worlds, World(line))
	}
	return PossibleWorlds{Worlds: worlds, key: key}
}

func BuildPossibleWorldsDeprecated(worlds []World) PossibleWorlds {
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
