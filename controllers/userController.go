package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	dababase "github.com/JacoboRodicio/daily-dish-server/database"
	"github.com/JacoboRodicio/daily-dish-server/helpers"
	"github.com/JacoboRodicio/daily-dish-server/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChangePasswordStruct struct {
	Oldpassword *string
	Newpassword *string
}

var userCollection *mongo.Collection = dababase.OpenCollection(dababase.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("The email or the password is incorrect")
		check = false
	}
	return check, msg
}

// TODO: Check validation
// Possible Issue: Creating an user even when the validation fails
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println()
		fmt.Printf("User: %v", user)
		fmt.Println()
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		count, emailErr := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if emailErr != nil {
			log.Panic(emailErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error ocurred while checking for the email"})
		}
		password := HashPassword(*user.Password)
		user.Password = &password

		count, phoneErr := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if phoneErr != nil {
			log.Panic(phoneErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error ocurred while checking for the phone"})
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This email or phone number already exists"})
		}
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, *&user.UserID)
		user.Token = &token
		user.RefreshToken = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user, foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := userCollection.FindOne(ctx, gin.H{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or Password is incorrect"})
			return
		}
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error in password check": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		}
		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, *&foundUser.UserID)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.UserID)
		err = userCollection.FindOne(ctx, bson.M{"userid": foundUser.UserID}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}

}

// Only available to Administrators
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, errPage := strconv.Atoi(c.Query("page"))
		if errPage != nil || page < 1 {
			page = 1
		}
		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		groupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
				{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
			},
			},
		}
		projectState := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
			}},
		}
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			groupStage, projectState,
		})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error when returning the list of users"})
			return
		}
		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}
		if len(allUsers) == 0 {
			c.JSON(http.StatusOK, []models.User{})
		} else {
			c.JSON(http.StatusOK, allUsers[0])
		}
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userid")
		if err := helpers.MatchUserTypeToUid(c, userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"userid": userID}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userid")
		docId, _ := primitive.ObjectIDFromHex(userId)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var updatedUser models.User
		var changePasswordStruct ChangePasswordStruct
		if err := c.BindJSON(&changePasswordStruct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			fmt.Println("err binding json: ", err)
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"_id": docId}).Decode(&updatedUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		passwordIsValid, msg := VerifyPassword(*changePasswordStruct.Oldpassword, *updatedUser.Password)
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		*updatedUser.Password = HashPassword(*changePasswordStruct.Newpassword)

		updatedTime, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		result, err := userCollection.ReplaceOne(
			ctx,
			bson.M{"_id": docId},
			bson.M{
				"firstname":    updatedUser.FirstName,
				"lastname":     updatedUser.LastName,
				"email":        updatedUser.Email,
				"phone":        updatedUser.Phone,
				"password":     updatedUser.Password,
				"updatedat":    updatedTime,
				"token":        updatedUser.Token,
				"refreshtoken": updatedUser.RefreshToken,
				"usertype":     updatedUser.UserType,
				"userid":       updatedUser.UserID,
				"createdat":    updatedUser.CreatedAt,
			})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			fmt.Println(err)
			return
		}
		c.JSON(http.StatusOK, result.ModifiedCount)
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userid")
		docId, _ := primitive.ObjectIDFromHex(userID)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user, updatedUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			fmt.Println("error updating user: ", err)
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"_id": docId}).Decode(&updatedUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if user.FirstName != nil {
			updatedUser.FirstName = user.FirstName
		}
		if user.LastName != nil {
			updatedUser.LastName = user.LastName
		}
		if user.Email != nil {
			updatedUser.Email = user.Email
		}
		if user.Phone != nil {
			updatedUser.Phone = user.Phone
		}
		if user.Password != nil {
			password := HashPassword(*user.Password)
			updatedUser.Password = &password
		}
		if user.UserType != nil {
			updatedUser.UserType = user.UserType
		}
		if user.UserID != "" {
			updatedUser.UserID = user.UserID
		}

		updatedTime, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		result, err := userCollection.ReplaceOne(
			ctx,
			bson.M{"_id": docId},
			bson.M{
				"firstname":    updatedUser.FirstName,
				"lastname":     updatedUser.LastName,
				"email":        updatedUser.Email,
				"phone":        updatedUser.Phone,
				"password":     updatedUser.Password,
				"updatedat":    updatedTime,
				"token":        updatedUser.Token,
				"refreshtoken": updatedUser.RefreshToken,
				"usertype":     updatedUser.UserType,
				"userid":       updatedUser.UserID,
				"createdat":    updatedUser.CreatedAt,
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			fmt.Println(err)
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result.ModifiedCount)
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userid")
		docId, _ := primitive.ObjectIDFromHex(userID)
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := userCollection.DeleteOne(ctx, bson.M{"_id": docId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result.DeletedCount)
	}
}
