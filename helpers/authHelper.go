package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// TODO: Is is a named return? i.e. can we save the lines 14 & 16? CHECK!
func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("usertype")
	err = nil
	if userType != role {
		err = errors.New("Unauthorized to access this resource")
		return err
	}
	return err
}

// TODO: Is is a named return? i.e. can we save the line 29? CHECK!
func MatchUserTypeToUid(c *gin.Context, userID string) (err error) {
	userType := c.GetString("usertype")
	uid := c.GetString("uid")
	err = nil

	if userType == "USER" && uid != userID {
		err = errors.New("Unauthorized to access this resource")
	}
	CheckUserType(c, userType)
	return err
}
