package models

import (
  "go.mongodb.org/mongo-driver/bson/primitive"
)


type Tag struct {
	ID              primitive.ObjectID `bson:"id"`
	Name            *string            `json:"name"`
}
