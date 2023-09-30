package models

type Category struct {
	Name        *string            `bson:"name"`
	Description *string            `bson:"description"`
	Favourite   *bool              `bson:"favourite"`
}
