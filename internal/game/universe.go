package game

import (
	"bytes"
	"encoding/binary"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/segmentio/fasthash/fnv1a"
)

type PossibleWorlds struct {
	Bits     roaring64.Bitmap
	count    int
	Capacity int
	hash     uint64
}

func (pw *PossibleWorlds) GetHash() uint64 {
	return pw.hash
}

func (pw *PossibleWorlds) SetHash() uint64 {
	bytes, _ := pw.Bits.ToBytes()
	pw.hash = fnv1a.Init64
	pw.hash = fnv1a.AddBytes64(pw.hash, bytes)
	return pw.hash
}

func (pw *PossibleWorlds) Add(index uint64) {
	pw.count += 1
	pw.Bits.Add(index)
}

func (pw *PossibleWorlds) Len() int {
	return pw.count
}

func (pw *PossibleWorlds) Indices() []uint64 {
	nums := pw.Bits.ToArray()
	return nums[:len(nums)-1]
}

func (pw *PossibleWorlds) ToBytes() []byte {
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.LittleEndian, pw.GetHash())
	binary.Write(&buffer, binary.LittleEndian, uint32(pw.Capacity))
	binary.Write(&buffer, binary.LittleEndian, uint32(pw.count))
	roaringBytes, _ := pw.Bits.ToBytes()
	buffer.Write(roaringBytes)
	return buffer.Bytes()
}

func PossibleWorldsFromBytes(data []byte) PossibleWorlds {
	buffer := bytes.NewBuffer(data[:16])
	var hash uint64
	var capacity uint32
	var count uint32
	binary.Read(buffer, binary.LittleEndian, &hash)
	binary.Read(buffer, binary.LittleEndian, &capacity)
	binary.Read(buffer, binary.LittleEndian, &count)

	bitmap := *roaring64.BitmapOf(uint64(capacity) + 1)
	bitmap.UnmarshalBinary(data[16:])

	return PossibleWorlds{
		hash:     hash,
		Capacity: int(capacity),
		count:    int(count),
		Bits:     bitmap,
	}
}

func NewPossibleWorlds(length int) PossibleWorlds {
	length64 := uint64(length)
	pw := PossibleWorlds{}
	pw.Bits = *roaring64.BitmapOf(length64 + 1)
	pw.count = 0
	pw.Capacity = length
	return pw
}

type Universe struct {
	Worlds         []World
	WorldLength    int
	Actions        AvailableActions
	PossibleWorlds PossibleWorlds
	Size           int
}

var BUniverse = Universe{}

func InitUniverse(lines []string) Universe {
	universe := Universe{}
	universe.Size = len(lines)
	universe.Worlds = make([]World, universe.Size)
	universe.WorldLength = len(lines[0])
	pw := NewPossibleWorlds(universe.Size)
	for index, name := range lines {
		universe.Worlds[index] = NewWorld(name)
		pw.Add(uint64(index))
	}
	pw.SetHash()
	universe.PossibleWorlds = pw
	universe.Actions = buildAvailableActions(universe)
	return universe
}

func (u *Universe) GetWorldByIndex(index uint64) World {
	return u.Worlds[index]
}

type AvailableActions struct {
	actions      []Action
	worldActions []Action
	worldSize    int
}

func buildAvailableActions(universe Universe) AvailableActions {
	aa := AvailableActions{
		actions:      []Action{},
		worldActions: make([]Action, universe.Size),
		worldSize:    universe.WorldLength,
	}

	keys := make(map[int32]bool)
	for index, world := range universe.Worlds {
		worldAction := NewAction(world.String())
		aa.worldActions[index] = worldAction
		for _, runeValue := range world.String() {
			if _, ok := keys[runeValue]; !ok {
				keys[runeValue] = true
				aa.actions = append(aa.actions, NewAction(string(runeValue)))
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
