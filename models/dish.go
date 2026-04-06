package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Dish struct {
	Name                  *string               `json:"name"`
	Description           *string               `json:"description"`
	Fat                   *float64              `json:"fat"`
	Ingredients           *[]primitive.ObjectID `json:"ingredients" bson:"ingredients"`
	IngredientQuantities  *map[string]string    `json:"ingredientQuantities" bson:"ingredientQuantities"`
	Recipe                *string               `json:"recipe"`
	Calories              *float64              `json:"calories"`
	Protein               *float64              `json:"protein"`
	Carbs                 *float64              `json:"carbs"`
	PreparationTime       *string               `json:"preparationTime" bson:"preparationTime"`
	CookTime              *string               `json:"cooktime"`
	Tags                  *[]primitive.ObjectID `json:"tags"`
	Categories            *[]primitive.ObjectID `json:"categories" bson:"categories"`
	Public                *bool                 `json:"public"`
	UserID                *string               `json:"userid"`
	AudioURL              *string               `json:"audioUrl"`
}
