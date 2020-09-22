package game

import (
	"encoding/binary"
	"fmt"

	"github.com/Workiva/go-datastructures/bitarray"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/mitchellh/hashstructure"
)

type BPossibleWorlds struct {
	Bits bitarray.BitArray
}

func (pw *BPossibleWorlds) GetHash() []byte {
	nums := []byte{}
	iter := pw.Bits.Blocks()
	for iter.Next() {
		_, block := iter.Value()
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(block))
		nums = append(nums, b...)
	}
	return nums
}

func (pw *BPossibleWorlds) GetBit(index uint64) bool {
	bool, err := pw.Bits.GetBit(index)
	if err != nil {
		panic(err)
	}
	return bool
}

func (pw *BPossibleWorlds) SetBit(index uint64) {
	err := pw.Bits.SetBit(index)
	if err != nil {
		panic(err)
	}
}

func (pw *BPossibleWorlds) GetName(index uint64) string {
	return BUniverse.Worlds[index].Name
}

func (pw *BPossibleWorlds) ToNums() []uint64 {
	nums := pw.Bits.ToNums()
	return nums[:len(nums)-1]
}

func NewPossibleWorlds(length int) BPossibleWorlds {
	length64 := uint64(length)
	pw := BPossibleWorlds{}
	pw.Bits = bitarray.NewBitArray(length64 + 1)
	pw.Bits.Reset()
	pw.Bits.SetBit(length64)
	return pw
}

func BuildActionWorldHash(action Action, worlds BPossibleWorlds) string {
	actionHash, err := hashstructure.Hash(action, nil)
	if err != nil {
		panic(err)
	}
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, actionHash)
	h := worlds.GetHash()
	h = append(h, b...)
	return string(h)

}

func LogCacheWorlds(prefix string, hash string, action Action, response string, bpw BPossibleWorlds) {
	fmt.Printf(
		"%s | hash:%s, action:%s, response:%s, worlds:%+v ",
		prefix,
		hash,
		action,
		response,
		bpw,
	)
}

type BWorld struct {
	Index uint64
	Name  string
}

type Universe struct {
	Worlds  []BWorld
	Actions AvailableActions
}

var BUniverse = Universe{}
var BPossibleWorldsStore = make(map[string]BPossibleWorlds)
var BBPossibleWorldsStore = make(map[string]map[string]BPossibleWorlds)
var BCountsStore = make(map[string]utils.GroupedIncrementer)
var BRewardsStore = make(map[string]utils.GroupedIncrementer)

func InitUniverse(lines []string) BPossibleWorlds {
	universeLength := len(lines)
	BUniverse.Worlds = make([]BWorld, universeLength)
	pw := NewPossibleWorlds(universeLength)
	for index, name := range lines {
		BUniverse.Worlds[index] = BWorld{Index: uint64(index), Name: name}
		pw.SetBit(uint64(index))
	}
	worlds := BuildPossibleWorlds(lines)
	BUniverse.Actions = BuildAvailableActions(worlds)
	return pw
}
