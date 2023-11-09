package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"platform_api/collections"
	"platform_api/mq"
	"platform_api/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type ImageController struct {
	ImageService collections.ImageCollection
}

func NewImageController(client *mongo.Client) *ImageController {
	return &ImageController{ImageService: *collections.NewImageCollection(client)}
}

// var imageCollection *mongo.Collection = configs.OpenCollection(configs.Client, "image_builder")

// GetAllImages godoc
//
//	@Summary		Retrieve all images
//	@Description	Get all image records from the database
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.Image
//	@Failure		500	{object}	models.HTTPError
//	@Router			/image [get]
func (t ImageController) GetAllImages(c *gin.Context) {
	images, statusCode, err := t.ImageService.GetAllImages()
	if err != nil {
		handleError(
			c,
			statusCode,
			"Failed to retrieve images",
			err,
		)
	}

	c.JSON(statusCode, *images)
}

// GetImageByCorId godoc
//
//	@Summary		Retrieve an image by Correlation ID
//	@Description	Get a single image record by Correlation ID (corId)
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			corId	path		string	true	"Correlation ID"
//	@Success		200		{object}	models.Image
//	@Failure		400		{object}	models.HTTPError
//	@Failure		404		{object}	models.HTTPError
//	@Router			/image/{corId} [get]
func (t ImageController) GetImageByCorId(c *gin.Context) {
	corId := c.Param("corId")

	image, statusCode, err := t.ImageService.GetImageByCorId(corId)
	if err != nil {
		handleError(
			c,
			statusCode,
			"Failed to retrieve image",
			err,
		)
		return
	}

	c.JSON(statusCode, *image)
}

// GetImageByCreatorName godoc
//
//	@Summary		Retrieve images by creator's name
//	@Description	Get all image records from the database filtered by creator's name
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			creatorName	path		string	true	"Creator's Name"
//	@Success		200			{array}		models.Image
//	@Failure		400			{object}	models.HTTPError
//	@Failure		404			{object}	models.HTTPError
//	@Router			/image/byCreator/{creatorName} [get]
func (t ImageController) GetImageByCreatorName(c *gin.Context) {
	creatorName := c.Param("creatorName")

	images, statusCode, err := t.ImageService.GetImageByCreatorName(creatorName)
	if err != nil {
		handleError(
			c,
			statusCode,
			"Failed to retrieve images",
			err,
		)
		return
	}

	c.JSON(http.StatusOK, *images)
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

// UploadImage godoc
//
//	@Summary		Upload an image
//	@Description	Upload an image file and trigger image creation process
//	@Tags			images
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			imageName	formData	string					true	"Name of the Image"
//	@Param			creatorName	formData	string					true	"Name of the Creator"
//	@Param			imageTag	formData	string					true	"Tag of the Image"
//	@Param			imageFile	formData	file					true	"The image file to upload"
//	@Success		200			{object}	map[string]interface{}	"A map containing the correlation ID"
//	@Failure		400			{object}	models.HTTPError
//	@Failure		500			{object}	models.HTTPError
//	@Router			/image/upload [post]
func (t ImageController) UploadImage(c *gin.Context) {

	// create uploadImage message
	var req UploadImageMessage

	// parse the result
	if c.PostForm("imageName") == "" || c.PostForm("creatorName") == "" || c.PostForm("imageTag") == "" {
		handleError(
			c,
			http.StatusBadRequest,
			"Error",
			errors.New("image name, tag, and creatorName cannot be empty"),
		)
		return
	}

	// assign to Message
	req.ImageName = c.PostForm("imageName")
	req.ImageTag = c.PostForm("imageTag")
	req.CreatorName = c.PostForm("creatorName")

	statusCode, err := t.ImageService.CheckImageByImageAndCreatorName(req.ImageName, req.CreatorName)
	if err != nil {
		handleError(
			c,
			statusCode,
			"Error",
			err,
		)
		return
	}

	// get the formFile
	formFile, err := c.FormFile("imageFile")
	if err != nil {
		handleError(
			c,
			http.StatusInternalServerError,
			"Error",
			err,
		)
		return
	}

	// get the other keys
	imageName := c.PostForm("imageName")
	creatorName := c.PostForm("creatorName")
	fileName := formFile.Filename
	file, err := formFile.Open()
	if err != nil {
		handleError(
			c,
			http.StatusBadRequest,
			"Error",
			err,
		)
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
		handleError(
			c,
			http.StatusBadRequest,
			"Error",
			err,
		)
		return
	}

	// set eventStatus
	req.EventStatus = "imageCreating"

	// set corId
	req.CorID = corId

	log.Printf("Uploaded file to %s", req)

	// validate json before passing to mq
	valid := validator.New()
	err = valid.Struct(req)
	if err != nil {
		handleError(
			c,
			http.StatusBadRequest,
			"Error",
			err,
		)
		return
	}

	// marshall data
	jsonReq, err := json.Marshal(req)
	if err != nil {
		handleError(
			c,
			http.StatusBadRequest,
			"Error",
			errors.New("failed to unmarshal data"),
		)
		return
	}

	// publish to mq
	err = mq.Pub(mq.EXCHANGE_TOPIC_ROUTER, mq.ROUTE_IMAGE_BUILD, jsonReq)
	if err != nil {
		handleError(
			c,
			http.StatusBadRequest,
			"Error",
			errors.New("failed to format request"),
		)
		return
	}

	// response
	resp := map[string]interface{}{"corId": corId}
	if err != nil {
		handleError(
			c,
			http.StatusBadRequest,
			"Error",
			errors.New("failed to generate corId"),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}
