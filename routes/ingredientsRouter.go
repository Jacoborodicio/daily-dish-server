package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"
  "strconv"

	database "github.com/JacoboRodicio/daily-dish-server/database"
	"github.com/JacoboRodicio/daily-dish-server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
)

var ingredientCollection *mongo.Collection = database.OpenCollection(database.Client, "ingredients")

func AddIngredient(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var ingredient models.Ingredient
	if err := c.BindJSON(&ingredient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validateErr := validate.Struct(ingredient)
	if validateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validateErr.Error()})
		fmt.Println(validateErr)
		return
	}
	ingredient.ID = primitive.NewObjectID()
	result, insertErr := ingredientCollection.InsertOne(ctx, ingredient)
	if insertErr != nil {
		msg := fmt.Sprintf("The new ingredient was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result)
}
func GetIngredients(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
  page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
  limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "0"), 10, 64)
  offset := (page - 1) * limit
  
  options := options.Find()
  options.SetLimit(limit)
  options.SetSkip(offset)
  
	var ingredients []bson.M
	cursor, err := ingredientCollection.Find(ctx, bson.M{}, options)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if err = cursor.All(ctx, &ingredients); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println(ingredients)
	c.JSON(http.StatusOK, ingredients)
}
func GetingredientById(c *gin.Context) {
	ingredientID := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(ingredientID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var ingredient bson.M
	if err := ingredientCollection.FindOne(ctx, bson.M{"_id": docId}).Decode(&ingredient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println(ingredient)
	c.JSON(http.StatusOK, ingredient)
}

// Still needed to change the whole stuff related to the copied categories
func UpdateIngredient(c *gin.Context) {
	ingredientID := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(ingredientID)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var ingredient, updatedIngredient models.Ingredient
	if err := c.BindJSON(&ingredient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validateErr := validate.Struct(ingredient)
	if validateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validateErr.Error()})
		fmt.Println(validateErr)
		return
	}
	if err := ingredientCollection.FindOne(ctx, bson.M{"_id": docId}).Decode(&updatedIngredient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	// Just to update what we receive
	if ingredient.Name != nil {
		updatedIngredient.Name = ingredient.Name
	}
	if ingredient.Calories != nil {
		updatedIngredient.Calories = ingredient.Calories
	}
	if ingredient.Price != nil {
		updatedIngredient.Price = ingredient.Price
	}

	result, err := ingredientCollection.ReplaceOne(
		ctx,
		bson.M{"_id": docId},
		bson.M{
			"id":       updatedIngredient.ID,
			"name":     updatedIngredient.Name,
			"price":    updatedIngredient.Price,
			"calories": updatedIngredient.Calories,
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
func Deleteingredient(c *gin.Context) {
	ingredientId := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(ingredientId)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	result, err := ingredientCollection.DeleteOne(ctx, bson.M{"_id": docId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()
	c.JSON(http.StatusOK, result.DeletedCount)
}
