package routes

import (
	controller "github.com/JacoboRodicio/daily-dish-server/controllers"

	"github.com/JacoboRodicio/daily-dish-server/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET(api+"/users", controller.GetUsers())
	incomingRoutes.GET(api+"/users/:userid", controller.GetUser())
	incomingRoutes.PUT(api+"/users/:userid", controller.UpdateUser())
	incomingRoutes.DELETE(api+"/users/:userid", controller.DeleteUser())
	incomingRoutes.PUT(api+"/users/changepassword/:userid", controller.ChangePassword())

}
