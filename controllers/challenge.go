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

type ChallengeController struct{}
var challengeCollection *mongo.Collection = configs.OpenCollection(configs.Client, "challenge")


func (t ChallengeController) GetAllChallenges(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := options.Find()
	cursor, err := challengeCollection.Find(ctx, bson.M{}, options)
	if err != nil {
		panic(err)
	}

	defer cursor.Close(ctx)	

	var challenges []models.Challenge
	err = cursor.All(ctx, &challenges)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{Code: http.StatusInternalServerError, Message: "Failed to retrieve challenges"})
		return
	}

	c.JSON(http.StatusOK, challenges)
}


func (t ChallengeController) GetChallengeById(c *gin.Context) {
	challengeId := c.Param("id")
	if challengeId == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "ID cannot be empty"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var challenge models.Challenge

	filter := bson.D{{Key: "_id", Value: challengeId}}
	err := challengeCollection.FindOne(ctx, filter).Decode(&challenge)

	if err != nil {
		msg := "Failed to retrieve challenge"
		if err == mongo.ErrNoDocuments {
			msg = "No challenge found with given Challenge ID"
		}
		c.JSON(http.StatusNotFound, models.HTTPError{Code: http.StatusNotFound, Message: msg})
		return
	}

	c.JSON(http.StatusOK, challenge)
}

func (t ChallengeController) GetChallengeByEmail(c *gin.Context) {
	email := c.Param("email")	
	if email == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "Email cannot be empty"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: email}}
	cursor, err := challengeCollection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}

	defer cursor.Close(ctx)

	var challenges []models.Challenge
	err = cursor.All(ctx, &challenges)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{Code: http.StatusInternalServerError, Message: "Failed to retrieve challenges"})
		return
	}

	if len(challenges) == 0 {  
		c.JSON(http.StatusNotFound, models.HTTPError{Code: http.StatusNotFound, Message: "No challenges with email found"})
		return
	}

	c.JSON(http.StatusOK, challenges)
}