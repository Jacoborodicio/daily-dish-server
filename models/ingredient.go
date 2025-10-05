package models

type Ingredient struct {
	Name     *string  `json:"name"`
	Calories *float64 `json:"calories"`
	Proteins *float64 `json:"proteins"`
	Price    *float64 `json:"price"`
}
