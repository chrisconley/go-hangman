package game

import (
	"testing"

	"github.com/chrisconley/go-hangman/test"
)

func TestMarshallUint64(t *testing.T) {
	uint1Bytes := Uint64ToBytes(uint64(123))

	uint2 := BytesToUint64(uint1Bytes)

	test.AssertDeepEqual(uint2, uint64(123), t)
}
