package models

type Goal struct {
	Name         *string  `bson:"name"`
	Description  *string  `bson:"description"`
	Value        *float64 `bson:"value"`
	Unit         *string  `bson:"unit"`
	ParentGoalId *string  `bson:"parentGoalId"`
}
