package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/JacoboRodicio/daily-dish-server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection = OpenCollection(Client, "categories")

func AddCategory(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var category models.Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validateErr := validate.Struct(category)
	if validateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validateErr.Error()})
		fmt.Println(validateErr)
		return
	}
	category.ID = primitive.NewObjectID()
	result, insertErr := categoryCollection.InsertOne(ctx, category)
	if insertErr != nil {
		msg := fmt.Sprintf("The new category was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result)
}
func GetCategories(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var categories []bson.M
	cursor, err := categoryCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if err = cursor.All(ctx, &categories); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println(categories)
	c.JSON(http.StatusOK, categories)
}
func GetCategoryById(c *gin.Context) {
	categoryID := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(categoryID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var category bson.M
	if err := categoryCollection.FindOne(ctx, bson.M{"_id": docId}).Decode(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println(category)
	c.JSON(http.StatusOK, category)
}
func UpdateCategory(c *gin.Context) {
	categoryID := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(categoryID)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var category models.Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validateErr := validate.Struct(category)
	if validateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validateErr.Error()})
		fmt.Println(validateErr)
		return
	}
	result, err := categoryCollection.ReplaceOne(
		ctx,
		bson.M{"_id": docId},
		bson.M{
			"name":        category.Name,
			"description": category.Description,
			"favourite":   category.Favourite,
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
func DeleteCategory(c *gin.Context) {
	categoryId := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(categoryId)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	result, err := categoryCollection.DeleteOne(ctx, bson.M{"_id": docId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()
	c.JSON(http.StatusOK, result.DeletedCount)
}
