package controllers

import (
	"context"
	"net/http"
	"platform_api/configs"
	"platform_api/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChallengeBuilderController struct{}

var challengeBuilderCollection *mongo.Collection = configs.OpenCollection(configs.Client, "challenge_builder")

func (t ChallengeBuilderController) GetAllChallenges(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := options.Find()
	cursor, err := challengeBuilderCollection.Find(ctx, bson.M{}, options)
	if err != nil {
		panic(err)
	}

	defer cursor.Close(ctx)

	var challenges []models.ChallengeBuilder
	err = cursor.All(ctx, &challenges)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{Code: http.StatusInternalServerError, Message: "Failed to retrieve challenges"})
		return
	}

	c.JSON(http.StatusOK, challenges)
}

// retrieve a list of challenges by their corID
func (t ChallengeBuilderController) GetChallengeByCorID(c *gin.Context) {
	corId := c.Param("corId")
	if corId == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "corId cannot be empty"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var challenges models.ChallengeBuilder

	filter := bson.D{{Key: "_id", Value: corId}}
	err := challengeBuilderCollection.FindOne(ctx, filter).Decode(&challenges)

	if err != nil {
		msg := "Failed to retrieve challenge"
		if err == mongo.ErrNoDocuments {
			msg = "No challenge found with given challenge corId"
		}
		c.JSON(http.StatusNotFound, models.HTTPError{Code: http.StatusNotFound, Message: msg})
		return
	}

	c.JSON(http.StatusOK, challenges)
}
