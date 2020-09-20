package game

import (
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
	// test.AssertDeepEqual(
	// 	BuildActionWorldHash(action, pw1),
	// 	BuildActionWorldHash(action, pw2),
	// 	t,
	// )

}
