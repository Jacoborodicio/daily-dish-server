package models

type Dish struct {
	Name            *string            `json:"name"`
	Fat             *float64           `json:"fat"`
	Ingredients     *string            `json:"ingredients"`
	Recipe          *string            `json:"recipe"`
	Calories        *float64           `json:"calories"`
	PreparationTime *string            `json:"preparationTime"`
  Tags            *[]Tag             `json:tags` 
  // Category        *string            `json:"category"`
  // Image           *string            `json:"image"`
  // Steps           *string            `json:"steps"`
}
