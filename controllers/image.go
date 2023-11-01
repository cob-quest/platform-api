package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"platform_api/configs"
	"platform_api/models"
	"platform_api/mq"
	"platform_api/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ImageController struct{}

var imageCollection *mongo.Collection = configs.OpenCollection(configs.Client, "image_builder")

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

	filter := bson.D{{Key: "corId", Value: corId}}
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

// @Summary: Get a image by creatorName
func (t ImageController) GetImageByCreatorName(c *gin.Context) {
	creatorName := c.Param("creatorName")
	if creatorName == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "creatorName cannot be empty"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "creatorName", Value: creatorName}}
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
		c.JSON(http.StatusNotFound, models.HTTPError{Code: http.StatusNotFound, Message: "No images with creatorName found"})
		return
	}

	c.JSON(http.StatusOK, images)
}

// POST Handler Body
type UploadImageMessage struct {
	ImageName   string `json:"imageName" validate:"required"`
	CreatorName string `json:"creatorName" validate:"required"`
	ImageTag    string `json:"imageTag" validate:"required"`
	S3Path      string `json:"s3Path" validate:"required"`
	CorID       string `json:"corId" validate:"required"`
	EventStatus string `json:"eventStatus" validate:"required"`
}

// POST Handler to upload a image zip file
func (t ImageController) UploadImage(c *gin.Context) {

	// create uploadImage message
	var req UploadImageMessage


	// parse the result
	if c.PostForm("imageName") == "" || c.PostForm("creatorName") == "" || c.PostForm("imageTag") == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "Missing imageName, creatorName or imageTag"})
		return
	}

	// assign to Message
	req.ImageName = c.PostForm("imageName")
	req.ImageTag = c.PostForm("imageTag")
	req.CreatorName= c.PostForm("creatorName")


	// ctx for 10s
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// check if such imagename+tag exists under this creator, if exists return it exists
	var image models.Image
	filter := bson.D{{Key: "imageName", Value: req.ImageName}, {Key: "creatorName", Value: req.CreatorName}, {Key: "imageTag", Value: req.ImageTag}}
	
	err := imageCollection.FindOne(ctx, filter).Decode(&image)
	if err == nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "ImageName and imageTag exists"})
		return
	} 

	

	// get the formFile
	formFile, err := c.FormFile("imageFile")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	// get the other keys
	imageName := c.PostForm("imageName")
	creatorName := c.PostForm("creatorName")
	fileName := formFile.Filename
	file, err := formFile.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	// generate correlationId
	corId := uuid.New().String()
	log.Printf("Received values: %s, %s, %s and generated %s", imageName, creatorName, fileName, corId)

	// create image message
	req.S3Path = fmt.Sprintf("%s/%s-%s.zip", "challenge-zips", creatorName, corId)

	// upload file to s3 compatible object store
	uploader := services.GetUploader()
	err = uploader.UploadFile(file, req.S3Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "Failed to upload file"})
		return
	}

	// validate json before passing to mq
	valid := validator.New()
	err = valid.Struct(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "Failed to validate form"})
		return
	}

	// set eventStatus
	req.EventStatus = "imageCreating"

	// marshall data
	jsonReq, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "Failed to unmarshall data"})
		return
	}

	// publish to mq
	err = mq.Pub(mq.EXCHANGE_TOPIC_ROUTER, mq.ROUTE_IMAGE_BUILD, jsonReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "Failed to format request"})
		return
	}

	// response
	resp := map[string]interface{}{"corId": corId}
	if err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "Failed to generate corId"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
