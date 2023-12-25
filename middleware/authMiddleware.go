package middleware

import (
	"fmt"
	"net/http"

	helpers "github.com/JacoboRodicio/daily-dish-server/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("No Authorization Token provided")})
			c.Abort()
			return
		}
		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("firstname", claims.Firstname)
		c.Set("lastname", claims.Lastname)
		c.Set("usertype", claims.Usertype)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
