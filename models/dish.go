package models
import (
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type Dish struct {
	Name            *string            `json:"name"`
	Fat             *float64           `json:"fat"`
	Ingredients     *[]primitive.ObjectID            `bson:"ingredients"`
	Recipe          *string            `json:"recipe"`
	Calories        *float64           `json:"calories"`
	PreparationTime *string            `json:"preparationTime"`
  Tags            *[]Tag             `json:tags` 
  Categories        *[]primitive.ObjectID            `bson:"categories"`
  // Image           *string            `json:"image"`
  // Steps           *string            `json:"steps"`
}
