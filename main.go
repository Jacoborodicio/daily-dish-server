package main

import (
	"os"

	"github.com/JacoboRodicio/daily-dish-server/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	// Dishes endpoints
	router.POST("/dishes", routes.AddDish)
	router.GET("/dishes", routes.GetDishes)
	router.GET("/dishes/:id/", routes.GetDishById)
	router.PUT("/dishes/:id", routes.UpdateDish)
	router.DELETE("/dishes/:id", routes.DeleteDish)

	// Categories endpoints
	router.POST("/categories", routes.AddCategory)
	router.GET("/categories", routes.GetCategories)
	router.GET("/categories/:id/", routes.GetCategoryById)
	router.PUT("/categories/:id", routes.UpdateCategory)
	router.DELETE("/categories/:id", routes.DeleteCategory)

	// Ingredients endpoints
	router.POST("/ingredients", routes.AddIngredient)
	router.GET("/ingredients", routes.GetIngredients)
	router.GET("/ingredients/:id/", routes.GetingredientById)
	router.PUT("/ingredients/:id", routes.UpdateIngredient)
	router.DELETE("/ingredients/:id", routes.Deleteingredient)
	router.Run(":" + port)
}
