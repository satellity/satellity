package uuid

import (
	"log"

	_uuid "github.com/satori/go.uuid"
)

// Nil is special form of UUID that is specified to have all
// 128 bits set to zero.
var Nil = _uuid.Nil

// NewV4 is a custom UUID
func NewV4() _uuid.UUID {
	id, err := _uuid.NewV4()
	if err != nil {
		log.Panicln(err)
	}
	return id
}

// FromString parse a string to UUID
func FromString(id string) (_uuid.UUID, error) {
	return _uuid.FromString(id)
}

// FromBytes parse a []byte to UUID
func FromBytes(input []byte) (_uuid.UUID, error) {
	return _uuid.FromBytes(input)
}
