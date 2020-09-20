package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestPossibleWorldsGetHash(t *testing.T) {
	pw1 := NewPossibleWorlds(100)
	pw1.SetBit(1)
	pw1.SetBit(99)

	pw2 := NewPossibleWorlds(100)
	pw2.SetBit(2)
	pw2.SetBit(98)
	test.AssertDeepNotEqual(pw1.GetHash(), pw2.GetHash(), t)
}

func TestHashingConflicts(t *testing.T) {
	pw1 := NewPossibleWorlds(25011)
	pw1.SetBit(uint64(617))
	pw1.SetBit(uint64(618))

	pw2 := NewPossibleWorlds(25011)
	pw2.SetBit(uint64(809))
	pw2.SetBit(uint64(810))

	action := Action("t")

	test.AssertDeepNotEqual(
		BuildActionWorldHash(action, pw1),
		BuildActionWorldHash(action, pw2),
		t,
	)
}
