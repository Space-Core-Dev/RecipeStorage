package store

import (
	"encoding/binary"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoIdConverter[T bson.ObjectID] struct {
}

func (c MongoIdConverter[T]) Convert(id int64) T {
	var b [12]byte

	binary.BigEndian.PutUint64(b[4:], uint64(id))

	return T(bson.ObjectID(b))
}
