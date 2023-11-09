package controllers

import (
	"encoding/json"
	"net/http"
	"platform_api/collections"
	"platform_api/models"

	"platform_api/mq"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type AttemptController struct{
	AttemptCollection collections.AttemptCollection
}

func NewAttemptController(client *mongo.Client) *AttemptController {
	return &AttemptController{AttemptCollection: *collections.NewAttemptCollection(client)}
}

// var attemptCollection *mongo.Collection = configs.OpenCollection(configs.Client, "attempt")

func (t AttemptController) GetAllAttempt(c *gin.Context) {
	attempts, statusCode, err := t.AttemptCollection.GetAllAttempts()
	if err != nil {
		handleError(
			c,
			statusCode,
			"Failed to retrieve attempts",
			err,
		)
		return
	}

	c.JSON(statusCode, *attempts)
}

// GetOneAttemptByToken finds and returns a challenge attempt by its token.
//
//	@Summary		Retrieve attempt by token
//	@Description	Get details of a specific attempt by token
//	@Tags			attempt
//	@Produce		json
//	@Param			token	path		string			true	"Attempt Token"
//	@Success		200		{object}	models.Attempt	"Successfully retrieved the attempt"
//	@Failure		400		"Invalid token parameter"
//	@Failure		404		"Attempt not found"
//	@Failure		500		"Internal server error"
//	@Router			/platform/attempt/{token} [get]
func (t AttemptController) GetOneAttemptByToken(c *gin.Context) {
	token := c.Param("token")
	
	attemptSingle, statusCode, err := t.AttemptCollection.GetOneAttemptByToken(token)
	if err != nil {
		handleError(
			c,
			statusCode,
			"Failed to retrieve attempt",
			err,
		)
		return
	}

	c.JSON(statusCode, *attemptSingle)
}

// StartAttempt creates a new attempt for a challenge.
//
//	@Summary		Start a new challenge attempt
//	@Description	Begin a new attempt for a specified challenge
//	@Tags			attempt
//	@Accept			json
//	@Produce		json
//	@Param			AttemptBody	body		AttemptBody				true	"Start Attempt Request Body"
//	@Success		200			{object}	models.SuccessResponse	"Successfully started the attempt with corId"
//	@Failure		400			"Bad request when the body is not as per AttemptBody structure"
//	@Failure		500			"Internal server error"
//	@Router			/platform/attempt [post]
func (t AttemptController) StartAttempt(c *gin.Context) {

	// parse body
	var req models.AttemptBody
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

	attemptSingle, statusCode, err := t.AttemptCollection.GetOneAttemptByToken(req.Token)
	if err != nil {
		handleError(
			c,
			statusCode,
			"Failed to retrieve attempt",
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

func (t AttemptController) SubmitAttemptByToken(c *gin.Context) {
	var req models.AttemptSubmitBody
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

	statusCode, err := t.AttemptCollection.UpdateAnAttempt(&req, req.Token)
	if err != nil {
		handleError(
			c,
			statusCode,
			"Failed to update an attempt",
			err,
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"submitted": true},
	)
}

func (t AttemptController) GetAllAttemptsByParticipant(c *gin.Context) {

	participant := c.Param("participant")
	if participant == "" {
		handleError(
			c,
			http.StatusBadRequest,
			"Invalid token",
			errors.New("token cannot be empty"),
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "participant", Value: participant}}

	// Create an update document to update the value of the object.
	var attempts []models.Attempt
	cursor, err := attemptCollection.Find(ctx, filter)
	if err != nil {
		handleError(
			c,
			http.StatusBadRequest,
			"Request not found",
			err,
		)
	}
	err = cursor.All(ctx, &attempts)
	if err != nil {
		handleError(
			c,
			http.StatusBadRequest,
			"Request not found",
			err,
		)
	}

	c.JSON(
		http.StatusOK,
		attempts,
	)
}
