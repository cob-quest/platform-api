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
)

type AttemptController struct{}

// AttemptBody is the request body for an attempt
// @Description AttemptBody is used to validate the request body for starting or getting an attempt.
// @Name AttemptBody
type AttemptBody struct {
	Token string `json:"token" validate:"required"`
	// Email string `json:"email" validate:"required"`
	ChallengeName     string `json:"challengeName" bson:"challengeName"`
	CreatorName       string `json:"creatorName" bson:"creatorName"`
	Participant       string `json:"participant" bson:"participant"`
	ImageRegistryLink string `json:"imageRegistryLink" bson:"imageRegistryLink"`
	CorId             string `json:"corId"`
	EventStatus       string `json:"eventStatus"`
}

var attemptCollection *mongo.Collection = configs.OpenCollection(configs.Client, "attempt")

// @Summary Get attempt information
// @Description Retrieve information about a specific attempt
// @Tags attempt
// @Accept json
// @Produce json
// @Param body body AttemptBody true "AttemptBody"
// @Success 200 {object} AttemptModel
// @Failure 400 {object} ErrorResponse
// @Router /startchallenge [get]
func (t AttemptController) GetAttempt(c *gin.Context) {
	var req AttemptBody
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		handleError(
			c,
			http.StatusInternalServerError,
			"Invalid token",
			errors.New("token cannot be empty"),
		)
		return
	}

	// ctx for 10s
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find the challenge with the specified ID
	var attemptSingle models.Attempt
	filter := bson.D{
		{Key: "token", Value: token},
		// {Key: "email", Value: req.Email},
	}
	err = attemptCollection.FindOne(ctx, filter).Decode(&attemptSingle)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		handleError(
			c,
			http.StatusInternalServerError,
			"Invalid Token or Email",
			err,
		)
		return
	} else if err != nil {
		handleError(
			c,
			http.StatusInternalServerError,
			"Invalid Token or Email",
			err,
		)
		return
	}

	c.JSON(http.StatusOK, attemptSingle)
}

// @Summary Start a new challenge attempt
// @Description Begin a new attempt for a challenge
// @Tags attempt
// @Accept json
// @Produce json
// @Param body body AttemptBody true "AttemptBody"
// @Success 200 {object} AttemptStartResponse
// @Failure 400 {object} ErrorResponse
// @Router /startchallenge [post]
func (t AttemptController) StartAttempt(c *gin.Context) {

	// parse body
	var req AttemptBody
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

	// validate json
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

	// find the challenge with the specified ID
	var attemptSingle models.Attempt
	filter := bson.D{
		{Key: "token", Value: req.Token},
		// {Key: "email", Value: req.Email},
	}
	err = attemptCollection.FindOne(ctx, filter).Decode(&attemptSingle)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		handleError(
			c,
			http.StatusBadRequest,
			"Invalid token",
			err,
		)
		return
	} else if err != nil {
		handleError(
			c,
			http.StatusBadRequest,
			"Invalid Token",
			err,
		)
		return
	}

	req.CorId = uuid.NewString()

	// marshall data for queue
	data, err := json.Marshal(attemptSingle)
	if err != nil {
		handleError(
			c,
			http.StatusInternalServerError,
			"Failed to marshall attempt to JSON",
			err,
		)
		return
	}

	err = json.Unmarshal(data, &req)
	if err != nil {
		handleError(
			c,
			http.StatusInternalServerError,
			"Failed to unmarshall attempt JSON to attempt Body",
			err,
		)
		return
	}

	// set eventStatus
	req.EventStatus = "challengeStarting"

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
	err = mq.Pub(mq.EXCHANGE_TOPIC_ROUTER, mq.ROUTE_CHALLENGE_START, jsonReq)
	if err != nil {
		handleError(
			c,
			http.StatusInternalServerError,
			"Failed to publish message",
			err,
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"corId": req.CorId},
	)
}

// POST Handler Body
type AttemptSubmitBody struct {
	Token string `json:"token" validate:"required"`
	Result float64 `json:"result" validate:"required"`
}

func (t AttemptController) SubmitAttemptByToken(c *gin.Context) {
	var req AttemptSubmitBody
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

	// validate json
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "token", Value: req.Token}}

	// Create an update document to update the value of the object.
	var attempt models.Attempt
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "result", Value: req.Result}}}}
	err = attemptCollection.FindOneAndUpdate(ctx, filter, update).Decode(&attempt)

	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		handleError(
			c,
			http.StatusBadRequest,
			"Invalid token",
			err,
		)
		return
	} else if err != nil {
		handleError(
			c,
			http.StatusBadRequest,
			"Invalid Token",
			err,
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"submitted": true},
	)
}
