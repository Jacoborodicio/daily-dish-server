package routes

import (
	controller "github.com/JacoboRodicio/daily-dish-server/controllers"

	"github.com/gin-gonic/gin"
)

const api = "/daily-dishes/api/v1/"

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST(api+"users/signup", controller.Signup())
	incomingRoutes.POST(api+"users/login", controller.Login())
}
