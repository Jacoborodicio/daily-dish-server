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

	const api = "/daily-dishes/api/v1/"
	// Dishes endpoints
	router.POST(api+"/dishes", routes.AddDish)
	router.POST(api+"/dishes/:id/audio", routes.HandleAudioUpload)
	router.GET(api+"/dishes/:id/audio/:fileName", routes.HandleAudioDownload)
	router.DELETE(api+"/dishes/:id/audio", routes.DeleteAudio)
	router.GET(api+"/dishes", routes.GetDishes)
	router.GET(api+"/dishes/:id", routes.GetDishById)
	router.PUT(api+"/dishes/:id", routes.UpdateDish)
	router.DELETE(api+"/dishes/:id", routes.DeleteDish)

	// Categories endpoints
	router.POST(api+"/categories", routes.AddCategory)
	router.GET(api+"/categories", routes.GetCategories)
	router.GET(api+"/categories/:id", routes.GetCategoryById)
	router.PUT(api+"/categories/:id", routes.UpdateCategory)
	router.DELETE(api+"/categories/:id", routes.DeleteCategory)

	// Ingredients endpoints
	router.POST(api+"/ingredients", routes.AddIngredient)
	router.GET(api+"/ingredients", routes.GetIngredients)
	router.GET(api+"/ingredients/:id", routes.GetIngredientById)
	router.PUT(api+"/ingredients/:id", routes.UpdateIngredient)
	router.DELETE(api+"/ingredients/:id", routes.DeleteIngredient)

	// Tags endpoints
	router.POST(api+"/tags", routes.AddTag)
	router.GET(api+"/tags", routes.GetTags)
	router.GET(api+"/tags/:id", routes.GetTagById)
	router.PUT(api+"/tags/:id", routes.UpdateTag)
	router.DELETE(api+"/tags/:id", routes.DeleteTag)

	router.Run(":" + port)
}
