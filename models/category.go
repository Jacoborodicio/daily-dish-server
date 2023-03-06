package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID          primitive.ObjectID `bson:"id"`
	Name        *string            `bson:"name"`
	Description *string            `bson:"description"`
	Favourite   *bool              `bson:"favourite"`
}
