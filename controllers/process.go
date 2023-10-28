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

type ProcessController struct{}

var processCollection *mongo.Collection = configs.OpenCollection(configs.Client, "process_engine")

func (t ProcessController) GetAllProcesses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := options.Find()
	cursor, err := processCollection.Find(ctx, bson.M{}, options)
	if err != nil {
		panic(err)
	}

	defer cursor.Close(ctx)

	var process []models.Process
	err = cursor.All(ctx, &process)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{Code: http.StatusInternalServerError, Message: "Failed to retrieve images"})
		return
	}

	c.JSON(http.StatusOK, process)
}

// retrieve a list of processes by their corID
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
	err := processCollection.FindOne(ctx, filter).Decode(&process)

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

func (t ProcessController) GetProcessByImageName(c *gin.Context) {
	imageName := c.Param("imageName")
	if imageName == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "imageName cannot be empty"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{
        {Key: "imageName", Value: imageName},
        {Key: "$text", Value: bson.D{{Key:"$search", Value: "image"}}},
    }
	cursor, err := processCollection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}

	defer cursor.Close(ctx)

	var process []models.Process
	err = cursor.All(ctx, &process)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{Code: http.StatusInternalServerError, Message: "Failed to retrieve processes. " + err.Error()})
		return
	}

	if len(process) == 0 {
		c.JSON(http.StatusNotFound, models.HTTPError{Code: http.StatusNotFound, Message: "No image process found with imageName"})
		return
	}

	c.JSON(http.StatusOK, process)
}
