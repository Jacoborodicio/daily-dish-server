package models

type Ingredient struct {
	Name     *string            `json:"name"`
	Calories *float64           `json:"calories"`
	Price    *float64           `json:"price"`
}
