package controllers

import (
	"encoding/json"
	"net/http"
	"platform_api/collections"
	"platform_api/mq"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

// ChallengeController is the controller for handling challenges.
type ChallengeController struct {
	ChallengeCollection collections.ChallengeCollection
	ImageCollection     collections.ImageCollection
}

func NewChallengeController(client *mongo.Client) *ChallengeController {
	return &ChallengeController{
		ChallengeCollection: *collections.NewChallengeCollection(client),
		ImageCollection:     *collections.NewImageCollection(client),
	}
}

// @Summary		Get all challenges Aaaaaaaaaaa
// @Description	Retrieves a list of all challenges.
// @Tags			challenges
// @Produce		json
// @Success		200	{array}		models.Challenge
// @Failure		500	{object}	models.HTTPError	"Failed to retrieve challenges"
// @Router			/challenge [get]
func (t ChallengeController) GetAllChallenges(c *gin.Context) {
	challenges, statusCode, err := t.ChallengeCollection.GetAllChallenges()
	if err != nil {
		handleError(
			c,
			statusCode,
			"Failed to retrieve challenges",
			err,
		)
		return
	}

	c.JSON(statusCode, *challenges)
}

// @Summary		Get challenge by CorID
// @Description	Retrieves a challenge based on its CorID.
// @Tags			challenges
// @Produce		json
// @Param			corId	path		string	true	"CorID of the Challenge"
// @Success		200		{object}	models.Challenge
// @Failure		400		{object}	models.HTTPError	"Invalid corId"
// @Failure		404		{object}	models.HTTPError	"No challenge found with given corId"
// @Router			/challenge/{corId} [get]
func (t ChallengeController) GetChallengeByCorID(c *gin.Context) {
	corId := c.Param("corId")

	challenge, statusCode, err := t.ChallengeCollection.GetChallengeByCorID(corId)
	if err != nil {
		handleError(
			c,
			statusCode,
			"Failed to retrieve challenges",
			err,
		)
		return
	}

	c.JSON(statusCode, *challenge)
}

// @Summary		Get challenge by creator name
// @Description	Retrieves a list of challenges based on the creator's name.
// @Tags			challenges
// @Produce		json
// @Param			creatorName	path		string	true	"Name of the Challenge Creator"
// @Success		200			{array}		models.Challenge
// @Failure		400			{object}	models.HTTPError	"Invalid creatorName"
// @Failure		404			{object}	models.HTTPError	"No challenges with creatorName found"
// @Failure		500			{object}	models.HTTPError	"Failed to retrieve challenges"
// @Router			/challenge/creator/{creatorName} [get]
func (t ChallengeController) GetChallengeByCreatorName(c *gin.Context) {
	creatorName := c.Param("creatorName")

	challenges, statusCode, err := t.ChallengeCollection.GetChallengeByCreatorName(creatorName)
	if err != nil {
		handleError(
			c,
			statusCode,
			"Failed to retrieve challenges",
			err,
		)
		return
	}

	// success
	c.JSON(http.StatusOK, *challenges)
}

// CreateChallengeMessage is the expected body content for creating a challenge.
type CreateChallengeMessage struct {
	CorID         string   `json:"corId"`
	ImageName     string   `json:"imageName" validate:"required"`
	ImageTag      string   `json:"imageTag" validate:"required"`
	ChallengeName string   `json:"challengeName" validate:"required"`
	CreatorName   string   `json:"creatorName" validate:"required"`
	Duration      int      `json:"duration" validate:"required"`
	Participants  []string `json:"participants" validate:"required"`
	EventStatus   string   `json:"eventStatus"`
}

// @Summary		Create a new challenge
// @Description	Creates a new challenge with the provided details.
// @Tags			challenge
// @Accept			json
// @Produce		json
// @Param			challenge	body		CreateChallengeMessage	true	"Create Challenge Content"
// @Success		200			{object}	models.SuccessResponse
// @Failure		400			{object}	models.HTTPError	"Invalid request body"
// @Failure		400			{object}	models.HTTPError	"Challenge name already exists"
// @Failure		404			{object}	models.HTTPError	"No such image"
// @Failure		500			{object}	models.HTTPError	"Failed to marshal JSON or Failed to publish message"
// @Failure		500			{object}	models.HTTPError	"Error occured while retrieving image"
// @Router			/challenge [post]
func (t ChallengeController) CreateChallenge(c *gin.Context) {

	// create createChallenge message
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

	// check if image exists
	statusCode, err := t.ImageCollection.CheckImageExists(req.ImageName, req.ImageTag, req.CreatorName)
	if err != nil {
		handleError(
			c,
			statusCode,
			"Error",
			err,
		)
		return
	}

	// check if the challenge name already exists
	statusCode, err = t.ChallengeCollection.CheckChallengeByChallengeAndCreatorName(req.CreatorName, req.ChallengeName)
	if err != nil {
		handleError(
			c,
			statusCode,
			"Error",
			err,
		)
		return
	}

	// set uuid
	corId := uuid.New().String()
	req.CorID = corId

	// set eventStatus
	req.EventStatus = "challengeCreating"

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
	err = mq.Pub(mq.EXCHANGE_TOPIC_ROUTER, mq.ROUTE_CHALLENGE_CREATE, jsonReq)
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
