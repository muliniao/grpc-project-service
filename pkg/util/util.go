package util

import (
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func WithID(id uuid.UUID) bson.M {
	return bson.M{"_id": id}
}

func WithUpdate(update bson.M) bson.M {
	return bson.M{"$set": update}
}
