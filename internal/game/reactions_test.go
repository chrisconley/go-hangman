package game

import (
	"fmt"
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

// func TestBigInt(t *testing.T) {
// 	// https://golang.org/pkg/math/big/
// 	// 4294967295             32 bit uint
// 	// 18446744073709551615   64 bit uint
// 	// 1099511627776          example Word value
// 	// Word is uint64
// 	// A BigInt of 10000 digits will be 157 Words (156.25*64 = 10000)
// 	// https://golang.org/ref/spec#Size_and_alignment_guarantees
// 	// type                                 size in bytes
// 	// byte, uint8, int8                     1
// 	// uint16, int16                         2
// 	// uint32, int32, float32                4
// 	// uint64, int64, float64, complex64     8
// 	// complex128                           16
// 	//
// 	// So 157 words takes up 157*8 = 1256 bytes
// 	//
// 	// There is also these possibly relevant packages
// 	// https://golang.org/pkg/encoding/binary/
// 	// https://golang.org/pkg/bytes/
// 	// https://golang.org/pkg/math/bits/
// 	// Other bitarray implementations
// 	// https://github.com/damnever/bitarray - Mar 2017 - 10 stars
// 	//     Uses bytes & encoding/binary
// 	// https://godoc.org/github.com/Workiva/go-datastructures/bitarray - Mar 2020 - 5.7k
// 	//     Looks like it just uses uint64 directly
// 	//     Uses bytes & encoding/binary in its encoding file
// 	//     Has lots of other data structures
// 	// https://github.com/golang-collections/go-datastructures - Feb 2015
// 	//     Forked from Workiva/go-datastructures
// 	var z1 big.Int
// 	// z1.Exp(big.NewInt(10), big.NewInt(99), nil)
// 	z1.SetBit(&z1, 1000, 1)
// 	fmt.Printf("word length %d", len(z1.Bits()))
// 	test.AssertDeepEqual(z1.Bits(), 12, t)
// }

// func TestBitArrays(t *testing.T) {
// 	lines := []string{"abc", "baa", "cab"}
// 	universeLength := uint64(len(lines))
// 	BUniverse = BUniverse[:universeLength]
// 	pw := NewPossibleWorlds(universeLength)
// 	for index, name := range lines {
// 		BUniverse[index] = BWorld{Index: uint64(index), Name: name}
// 		pw.Bits.SetBit(uint64(index))
// 	}
// 	BPossibleWorldsStore[pw.GetHash()] = pw
// 	test.AssertDeepEqual(pw.GetBit(0), true, t)
// 	test.AssertDeepEqual(pw.GetBit(1), true, t)
// 	test.AssertDeepEqual(pw.GetBit(2), true, t)
// 	fmt.Printf("pw: %v with hash: %v\n", pw.Bits.ToNums(), pw.GetHash())
// 	fmt.Printf("pw: %v\n", BPossibleWorldsStore)
// 	fmt.Printf("pw: %v\n", BPossibleWorldsStore[pw.GetHash()].Bits.ToNums())
// }

func TestReactionExplorerGet(t *testing.T) {
	re := NewReactionExplorer()
	action := Action("a")
	response := Response("a--")
	actual := re.Get(action, response)
	test.AssertDeepEqual(actual, make([]World, 0), t)

	re.Add(action, response, 0.75, "abc")
	actual = re.Get(action, response)
	test.AssertDeepEqual(actual, []World{"abc"}, t)

	re.Add(action, response, 0.5, "aaa")
	actual = re.Get(action, response)
	test.AssertDeepEqual(actual, []World{"abc", "aaa"}, t)
}

func TestReactionExplorerActionCounts(t *testing.T) {
	re := NewReactionExplorer()
	action := Action("a")

	actual := re.ActionCounts(action)
	test.AssertDeepEqual(actual, []float64{}, t)

	re.Add(action, Response("a--"), 0.75, "abc")
	re.Add(action, Response("a--"), 1.0, "aaa")
	re.Add(action, Response("!"), 0.0, "bbc")

	actual = re.ActionCounts(action)
	test.AssertDeepEqual(actual, []float64{2.0, 1.0}, t)
}

func TestReactionExplorerActionTotal(t *testing.T) {
	re := NewReactionExplorer()
	action := Action("a")

	actual := re.ActionTotal(action)
	test.AssertDeepEqual(actual, 0.0, t)

	re.Add(action, Response("a--"), 0.75, "abc")
	re.Add(action, Response("a--"), 1.0, "aaa")
	re.Add(action, Response("!"), 0.0, "bbc")

	actual = re.ActionTotal(action)
	test.AssertDeepEqual(actual, 3.0, t)
}

func TestReactionExplorerActionRewardTotal(t *testing.T) {
	re := NewReactionExplorer()
	action := Action("a")

	actual := re.ActionRewardTotal(action)
	test.AssertDeepEqual(actual, 0.0, t)

	re.Add(action, Response("a--"), 0.75, "abc")
	re.Add(action, Response("a--"), 1.0, "aaa")
	re.Add(action, Response("!"), 0.2, "bbc")

	actual = re.ActionRewardTotal(action)
	test.AssertDeepEqual(actual, 1.95, t)
}

func TestReactionExplorerActions(t *testing.T) {
	re := NewReactionExplorer()
	test.AssertDeepEqual(re.Actions(), []Action{}, t)

	re.Add(Action("a"), Response("a--"), 0.75, "abc")
	re.Add(Action("a"), Response("a--"), 1.0, "aaa")
	re.Add(Action("b"), Response("!"), 0.0, "aac")

	test.AssertDeepEqual(re.Actions(), []Action{"a", "b"}, t)
}

func BenchmarkActionsSmall(b *testing.B) {
	re := NewReactionExplorer()
	numElements := 10
	for i := 0; i < numElements; i++ {
		re.Add(Action(fmt.Sprintf("%d", i)), "!", 0.0, "a")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		re.Actions()
	}
}

func BenchmarkActionsLarge(b *testing.B) {
	re := NewReactionExplorer()
	numElements := 10000
	for i := 0; i < numElements; i++ {
		re.Add(Action(fmt.Sprintf("%d", i)), "!", 0.0, "a")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		re.Actions()
	}
}

func TestPossibleWorldsGetHash(t *testing.T) {
	pw1 := NewPossibleWorlds(100)
	pw1.Bits.SetBit(1)
	pw1.Bits.SetBit(99)

	pw2 := NewPossibleWorlds(100)
	pw2.Bits.SetBit(2)
	pw2.Bits.SetBit(98)
	test.AssertDeepNotEqual(pw1.GetHash(), pw2.GetHash(), t)
}
