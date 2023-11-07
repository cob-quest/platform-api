package controllers

import (
	"context"
	"errors"
	"net/http"
	"platform_api/configs"
	"platform_api/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProcessController struct{}

var processCollection *mongo.Collection = configs.OpenCollection(configs.Client, "process_engine")

// GetAllProcesses godoc
//	@Summary		Retrieves all processes
//	@Description	Get all the processes from the process engine
//	@Tags			processes
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.Process
//	@Failure		500	{object}	models.HTTPError
//	@Router			/processes [get]
func (t ProcessController) GetAllProcesses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := options.Find()
	cursor, err := processCollection.Find(ctx, bson.M{}, options)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)

	// retrieve the process
	var process []models.Process
	err = cursor.All(ctx, &process)
	if err != nil {
		handleError(
			c,
			http.StatusInternalServerError,
			"Failed to retrieve processes",
			err,
		)
		return
	}

	c.JSON(http.StatusOK, process)
}

// GetProcessByCorID godoc
//	@Summary		Retrieves a process by Correlation ID
//	@Description	Retrieve a list of processes by their Correlation ID
//	@Tags			processes
//	@Accept			json
//	@Produce		json
//	@Param			corId	path		string	true	"Correlation ID"
//	@Success		200		{object}	models.Process
//	@Failure		400		{object}	models.HTTPError
//	@Failure		404		{object}	models.HTTPError
//	@Router			/processes/{corId} [get]
func (t ProcessController) GetProcessByCorID(c *gin.Context) {
	corId := c.Param("corId")
	if corId == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "corId cannot be empty"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var process models.Process

	filter := bson.D{{Key: "corId", Value: corId}}
	opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}}) // -1 for descending order
	err := processCollection.FindOne(ctx, filter, opts).Decode(&process)

	if err != nil {
		msg := "Failed to retrieve challenge"
		if err == mongo.ErrNoDocuments {
			msg = "No process found with given process corId"
		}
		c.JSON(http.StatusNotFound, models.HTTPError{Code: http.StatusNotFound, Message: msg})
		return
	}

	c.JSON(http.StatusOK, process)
}

// GetProcessByCreatorName godoc
//	@Summary		Retrieves processes by Creator Name
//	@Description	Retrieve a list of processes filtered by Creator Name
//	@Tags			processes
//	@Accept			json
//	@Produce		json
//	@Param			creatorName	path		string	true	"Creator's Name"
//	@Success		200			{array}		models.Process
//	@Failure		400			{object}	models.HTTPError
//	@Failure		404			{object}	models.HTTPError
//	@Router			/processes/byCreator/{creatorName} [get]
func (t ProcessController) GetProcessByCreatorName(c *gin.Context) {
	creatorName := c.Param("creatorName")
	if creatorName == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "creatorName cannot be empty"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "creatorName", Value: creatorName}}
	cursor, err := processCollection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}

	defer cursor.Close(ctx)

	var process []models.Process
	err = cursor.All(ctx, &process)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{Code: http.StatusInternalServerError, Message: "Failed to retrieve processes"})
		return
	}

	if len(process) == 0 {
		c.JSON(http.StatusNotFound, models.HTTPError{Code: http.StatusNotFound, Message: "No process with creatorName found"})
		return
	}

	c.JSON(http.StatusOK, process)
}

// GetProcessStatusByCorId godoc
//	@Summary		Retrieves the status of a process by Correlation ID
//	@Description	Get the most recent status of a specific process by Correlation ID
//	@Tags			processes
//	@Accept			json
//	@Produce		json
//	@Param			corId	path		string	true	"Correlation ID"
//	@Success		200		{object}	models.Process
//	@Failure		400		{object}	models.HTTPError
//	@Failure		404		{object}	models.HTTPError
//	@Failure		500		{object}	models.HTTPError
//	@Router			/processes/status/{corId} [get]
func (t ProcessController) GetProcessStatusByCorId(c *gin.Context) {
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

    var process models.Process
    filter := bson.D{
        {Key: "corId", Value: corId},
    }

    options := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})
    err := processCollection.FindOne(ctx, filter, options).Decode(&process)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		handleError(
			c,
			http.StatusNotFound,
			"Invalid CorId given",
			err,
		)
		return
	} else if err != nil {
		handleError(
			c,
			http.StatusInternalServerError,
			"Error",
			err,
		)
		return
	}
    c.JSON(http.StatusOK, process)
}