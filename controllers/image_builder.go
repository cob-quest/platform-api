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

type ImageBuilderController struct{}

var imageBuilderCollection *mongo.Collection = configs.OpenCollection(configs.Client, "image_builder")

func (t ImageBuilderController) GetAllImages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := options.Find()
	cursor, err := imageBuilderCollection.Find(ctx, bson.M{}, options)
	if err != nil {
		panic(err)
	}

	defer cursor.Close(ctx)

	var images []models.ImageBuilder
	err = cursor.All(ctx, &images)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{Code: http.StatusInternalServerError, Message: "Failed to retrieve images"})
		return
	}

	c.JSON(http.StatusOK, images)
}

// retrieve a list of images by their corID
func (t ImageBuilderController) GetImageByCorID(c *gin.Context) {
	corId := c.Param("corId")
	if corId == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "corId cannot be empty"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var images models.ImageBuilder

	filter := bson.D{{Key: "_id", Value: corId}}
	err := imageBuilderCollection.FindOne(ctx, filter).Decode(&images)

	if err != nil {
		msg := "Failed to retrieve image"
		if err == mongo.ErrNoDocuments {
			msg = "No challenge found with given image corId"
		}
		c.JSON(http.StatusNotFound, models.HTTPError{Code: http.StatusNotFound, Message: msg})
		return
	}

	c.JSON(http.StatusOK, images)
}
