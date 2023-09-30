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

var tagCollection *mongo.Collection = database.OpenCollection(database.Client, "tags")

func AddTag(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var tag models.Tag
	if err := c.BindJSON(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validateErr := validate.Struct(tag)
if validateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validateErr.Error()})
		fmt.Println(validateErr)
		return
	}
	tag.ID = primitive.NewObjectID()
	result, insertErr := tagCollection.InsertOne(ctx, tag)
	if insertErr != nil {
		msg := fmt.Sprintf("The new tag was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result)
}
func GetTags(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
  page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
  limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "0"), 10, 64)
  offset := (page - 1) * limit
  
  options := options.Find()
  options.SetLimit(limit)
  options.SetSkip(offset)
  
	var tags []bson.M
	cursor, err := tagCollection.Find(ctx, bson.M{}, options)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if err = cursor.All(ctx, &tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println(tags)
	c.JSON(http.StatusOK, tags)
}
func GetTagById(c *gin.Context) {
	tagID := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(tagID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var tag bson.M
	if err := tagCollection.FindOne(ctx, bson.M{"_id": docId}).Decode(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println(tag)
	c.JSON(http.StatusOK, tag)
}
func UpdateTag(c *gin.Context) {
	tagID := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(tagID)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var tag, updatedTag models.Tag
	if err := c.BindJSON(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validateErr := validate.Struct(tag)
	if validateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validateErr.Error()})
		fmt.Println(validateErr)
		return
	}
	if err := tagCollection.FindOne(ctx, bson.M{"_id": docId}).Decode(&updatedTag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	// Just to update what we receive
	if tag.Name != nil {
		updatedTag.Name = tag.Name
	}
	result, err := tagCollection.ReplaceOne(
		ctx,
		bson.M{"_id": docId},
		bson.M{
			"id":          updatedTag.ID,
			"name":        updatedTag.Name,
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
func DeleteTag(c *gin.Context) {
	tagId := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(tagId)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	result, err := tagCollection.DeleteOne(ctx, bson.M{"_id": docId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()
	c.JSON(http.StatusOK, result.DeletedCount)
}
