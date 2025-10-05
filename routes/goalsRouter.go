package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	database "github.com/JacoboRodicio/daily-dish-server/database"
	"github.com/JacoboRodicio/daily-dish-server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var goalsCollection *mongo.Collection = database.OpenCollection(database.Client, "goals")

func AddGoal(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var goal models.Goal
	if err := c.BindJSON(&goal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validateErr := validate.Struct(goal)
	if validateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validateErr.Error()})
		fmt.Println(validateErr)
		return
	}
	result, insertErr := goalsCollection.InsertOne(ctx, goal)
	if insertErr != nil {
		msg := fmt.Sprintf("The new goal was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}
	c.JSON(http.StatusOK, result)
}
func GetGoals(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "0"), 10, 64)
	offset := (page - 1) * limit

	options := options.Find()
	options.SetLimit(limit)
	options.SetSkip(offset)

	var goals []bson.M
	cursor, err := goalsCollection.Find(ctx, bson.M{}, options)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if err = cursor.All(ctx, &goals); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	fmt.Println(goals)
	c.JSON(http.StatusOK, goals)
}
func GetGoalById(c *gin.Context) {
	goalID := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(goalID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var goal bson.M
	if err := tagCollection.FindOne(ctx, bson.M{"_id": docId}).Decode(&goal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	fmt.Println(goal)
	c.JSON(http.StatusOK, goal)
}
func UpdateGoal(c *gin.Context) {
	goalID := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(goalID)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var goal, updatedGoal models.Goal
	if err := c.BindJSON(&goal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validateErr := validate.Struct(goal)
	if validateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validateErr.Error()})
		fmt.Println(validateErr)
		return
	}
	if err := tagCollection.FindOne(ctx, bson.M{"_id": docId}).Decode(&updatedGoal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	// TODO: It needs to be updated to allow the update of other fields
	// Just to update what we receive
	if goal.Name != nil {
		updatedGoal.Name = goal.Name
	}
	result, err := goalsCollection.ReplaceOne(
		ctx,
		bson.M{"_id": docId},
		bson.M{
			"name": updatedGoal.Name,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, result.ModifiedCount)
}
func DeleteGoal(c *gin.Context) {
	goalId := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(goalId)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	result, err := goalsCollection.DeleteOne(ctx, bson.M{"_id": docId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, result.DeletedCount)
}
