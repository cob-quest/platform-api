package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"platform_api/configs"
	"platform_api/models"
	"platform_api/mq"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChallengeController struct{}

var challengeCollection *mongo.Collection = configs.OpenCollection(configs.Client, "challenge")

// Gets all the challenges in MongoDB
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
		handleError(
			c,
			http.StatusInternalServerError,
			"Failed to retrieve challenges",
			err,
		)
		return
	}

	c.JSON(http.StatusOK, challenges)
}

// Gets all Challenges by the given corId
func (t ChallengeController) GetChallengeByCorID(c *gin.Context) {
	corId := c.Param("corId")
	if corId == "" {
		handleError(
			c,
			http.StatusInternalServerError,
			"Invalid corId",
			errors.New("corId cannot be empty"),
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var challenges models.Challenge

	filter := bson.D{{Key: "corId", Value: corId}}
	err := challengeCollection.FindOne(ctx, filter).Decode(&challenges)
	if err != nil {
		msg := "Failed to retrieve challenge"
		if errors.Is(err, mongo.ErrNoDocuments) {
			msg = "No challenge found with given challenge corId"
		}
		handleError(
			c,
			http.StatusBadRequest,
			msg,
			err,
		)
		return
	}

	c.JSON(http.StatusOK, challenges)
}

// Gets Challenge by a Creator Name
func (t ChallengeController) GetChallengeByCreatorName(c *gin.Context) {

	creatorName := c.Param("creatorName")
	if creatorName == "" {
		handleError(
			c,
			http.StatusBadRequest,
			"Invalid creatorName",
			errors.New("creatorName cannot be empty"),
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// filter the challenges by the creatorName
	filter := bson.D{{Key: "creatorName", Value: creatorName}}
	cursor, err := challengeCollection.Find(ctx, filter)
	if err != nil {
		handleError(
			c,
			http.StatusBadRequest,
			"Invalid request",
			err,
		)
		return
	}

	defer cursor.Close(ctx)

	var challenge []models.Challenge

	// get all the challenges
	err = cursor.All(ctx, &challenge)
	if err != nil {
		handleError(
			c,
			http.StatusInternalServerError,
			"Failed to retrieve challenges",
			err,
		)
		return
	}

	// if there are 0 challenges, respond with error
	if len(challenge) == 0 {
		handleError(
			c,
			http.StatusInternalServerError,
			"No challenges with creatorName found",
			err,
		)
		return
	}

	// success
	c.JSON(http.StatusOK, challenge)
}

// POST Handler Body
type CreateChallengeMessage struct {
	CorID         string   `json:"corId"`
	ImageName     string   `json:"imageName" validate:"required"`
	ChallengeName string   `json:"challengeName" validate:"required"`
	CreatorName   string   `json:"creatorName" validate:"required"`
	Duration      int      `json:"duration" validate:"required"`
	Participants  []string `json:"participants" validate:"required"`
}

// Creates a new challenge
func (t ChallengeController) CreateChallenge(c *gin.Context) {

	// create image message
	var req CreateChallengeMessage

	// parse the result
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		handleError(
			c,
			http.StatusBadRequest,
			"Invalid request body json",
			err,
		)
		return
	}

	// validate json before passing to mq
	v := validator.New()
	err = v.Struct(req)
	if err != nil {
		handleError(
			c,
			http.StatusBadRequest,
			"Invalid request body",
			err,
		)
		return
	}

	// ctx for 10s
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// check if the image name exists
	var image models.Image
	filter := bson.D{{Key: "imageName", Value: req.ImageName}}
	err = imageCollection.FindOne(ctx, filter).Decode(&image)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		handleError(
			c,
			http.StatusInternalServerError,
			"No such image",
			err,
		)
		return
	} else if err != nil {
		handleError(
			c,
			http.StatusInternalServerError,
			"Error occured while retrieving image",
			err,
		)
		return
	}

	// check if the challenge name already exists
	var challenge models.Challenge
	filter = bson.D{{Key: "challengeName", Value: req.ChallengeName}}
	err = challengeCollection.FindOne(ctx, filter).Decode(&challenge)
	// for the image to be valid, its name must not exist already
	if !(err != nil && errors.Is(err, mongo.ErrNoDocuments)) {
		handleError(
			c,
			http.StatusInternalServerError,
			"Challenge name already exists",
			err,
		)
		return
	}

	// set uuid
	corId := uuid.New().String()
	req.CorID = corId

	// marshall data
	jsonReq, err := json.Marshal(req)
	if err != nil {
		handleError(
			c,
			http.StatusInternalServerError,
			"Failed to marshall JSON",
			err,
		)
		return
	}

	// publish to mq
	err = mq.Pub(mq.EXCHANGE_TOPIC_TRIGGER, mq.ROUTE_CHALLENGE_CREATE, jsonReq)
	if err != nil {
		handleError(
			c,
			http.StatusInternalServerError,
			"Failed to publish message",
			err,
		)
		return
	}

	// response
	resp := map[string]interface{}{"corId": corId}
	c.JSON(http.StatusOK, resp)
}
