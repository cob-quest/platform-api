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

type ImageController struct{}

var imageCollection *mongo.Collection = configs.OpenCollection(configs.Client, "image")

func (t ImageController) GetAllImages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := options.Find()
	cursor, err := imageCollection.Find(ctx, bson.M{}, options)
	if err != nil {
		panic(err)
	}

	defer cursor.Close(ctx)

	var images []models.Image
	err = cursor.All(ctx, &images)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{Code: http.StatusInternalServerError, Message: "Failed to retrieve images"})
		return
	}

	c.JSON(http.StatusOK, images)
}

func (t ImageController) GetImageByCorId(c *gin.Context) {
	corId := c.Param("corId")
	if corId == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "ID cannot be empty"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var image models.Image

	filter := bson.D{{Key: "cor_id", Value: corId}}
	err := imageCollection.FindOne(ctx, filter).Decode(&image)

	if err != nil {
		msg := "Failed to retrieve image"
		if err == mongo.ErrNoDocuments {
			msg = "No image found with given Image ID"
		}
		c.JSON(http.StatusNotFound, models.HTTPError{Code: http.StatusNotFound, Message: msg})
		return
	}

	c.JSON(http.StatusOK, image)
}

// @Summary: Get a image by email
func (t ImageController) GetImageByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "Email cannot be empty"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: email}}
	cursor, err := imageCollection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}

	defer cursor.Close(ctx)

	var images []models.Image
	err = cursor.All(ctx, &images)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{Code: http.StatusInternalServerError, Message: "Failed to retrieve images"})
		return
	}

	if len(images) == 0 {
		c.JSON(http.StatusNotFound, models.HTTPError{Code: http.StatusNotFound, Message: "No images with email found"})
		return
	}

	c.JSON(http.StatusOK, images)
}
