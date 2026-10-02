package mongo

import (
	"encoding/binary"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IdConverter[T primitive.ObjectID] struct {
}

func (c IdConverter[T]) Convert(id int64) T {
	var b [12]byte

	binary.BigEndian.PutUint64(b[4:], uint64(id))

	return T(primitive.ObjectID(b))
}
