package storage

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSession struct {
	Id         primitive.ObjectID `bson:"_id"`
	Expiration time.Time          `bson:"expiration"`
}
type UserModel struct {
	Id       primitive.ObjectID `bson:"_id"`
	Login    string             `bson:"login"`
	Password []byte             `bson:"password"`
}
