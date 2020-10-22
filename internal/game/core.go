package game

import "github.com/segmentio/fasthash/fnv1a"

type World interface {
	HashedString
}

type world hashedString

func NewWorld(actionString string) World {
	return newHashedString(actionString)
}

type Action interface {
	HashedString
}

type action hashedString

func NewAction(actionString string) Action {
	return newHashedString(actionString)
}

func NewActionFromBytes(actionBytes []byte) Action {
	hs := hashedString{}
	hs.bytes = actionBytes
	hs.hash = fnv1a.HashBytes64(actionBytes)
	hs.length = len(actionBytes)
	return hs
}

type Response interface {
	HashedString
}

type response hashedString

func NewResponse(responseString string) Response {
	return newHashedString(responseString)
}

type HashedString interface {
	String() string
	Hash() uint64
	Bytes() []byte
	Matches(hs HashedString) bool
	Length() int
}

type hashedString struct {
	bytes  []byte
	hash   uint64
	length int
}

func newHashedString(str string) HashedString {
	hs := hashedString{}
	hs.bytes = []byte(str)
	hs.hash = fnv1a.HashBytes64(hs.bytes)
	hs.length = len(str)
	return hs
}

func (w hashedString) String() string {
	return string(w.bytes)
}

func (w hashedString) Hash() uint64 {
	return w.hash
}

func (w hashedString) Bytes() []byte {
	return w.bytes
}

func (w hashedString) Length() int {
	return w.length
}

func (w hashedString) Matches(other HashedString) bool {
	return w.Hash() == other.Hash()
}
