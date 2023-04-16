package routes

import (
	controller "github.com/JacoboRodicio/daily-dish-server/controllers"

	"github.com/JacoboRodicio/daily-dish-server/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("users/:user_id", controller.GetUser())
}