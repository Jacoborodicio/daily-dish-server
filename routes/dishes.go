package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/JacoboRodicio/daily-dish-server/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()
var dishCollection *mongo.Collection = OpenCollection(Client, "dishes")

func AddDish(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var dish models.Dish
	if err := c.BindJSON(&dish); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validateErr := validate.Struct(dish)
	if validateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validateErr.Error()})
		fmt.Println(validateErr)
		return
	}
	dish.ID = primitive.NewObjectID()
	result, insertErr := dishCollection.InsertOne(ctx, dish)
	if insertErr != nil {
		msg := fmt.Sprintf("The new dish was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result)
}
func GetDishes(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var dishes []bson.M
	cursor, err := dishCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if err = cursor.All(ctx, &dishes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println(dishes)
	c.JSON(http.StatusOK, dishes)
}
func GetDishById(c *gin.Context) {
	dishID := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(dishID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var dish bson.M
	if err := dishCollection.FindOne(ctx, bson.M{"_id": docId}).Decode(&dish); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println(dish)
	c.JSON(http.StatusOK, dish)
}

func UpdateDish(c *gin.Context) {
	dishId := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(dishId)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var dish models.Dish
	if err := c.BindJSON(&dish); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validateErr := validate.Struct(dish)
	if validateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validateErr.Error()})
		fmt.Println(validateErr)
		return
	}
	result, err := dishCollection.ReplaceOne(
		ctx,
		bson.M{"_id": docId},
		bson.M{
			"name":             dish.Name,
			"fat":              dish.Fat,
			"ingredients":      dish.Ingredients,
			"recipe":           dish.Recipe,
			"calories":         dish.Calories,
			"preparation-time": dish.PreparationTime,
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

func DeleteDish(c *gin.Context) {
	dishId := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(dishId)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	result, err := dishCollection.DeleteOne(ctx, bson.M{"_id": docId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()
	c.JSON(http.StatusOK, result.DeletedCount)
}
