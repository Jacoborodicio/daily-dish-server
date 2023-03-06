package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Dish struct {
	ID              primitive.ObjectID `bson:"id"`
	Name            *string            `json:"name"`
	Fat             *float64           `json:"fat"`
	Ingredients     *string            `json:"ingredients"`
	Recipe          *string            `json:"recipe"`
	Calories        *float64           `json:"calories"`
	PreparationTime *string            `json:"preparationTime"`
}
