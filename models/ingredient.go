package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ingredient struct {
	ID       primitive.ObjectID `bson:"id"`
	Name     *string            `json:"name"`
	Calories *float64           `json:"calories"`
	Price    *float64           `json:"price"`
}
