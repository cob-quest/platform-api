package controllers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"platform_api/configs"
	"platform_api/models"

	"github.com/gin-gonic/gin"

	// "go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/stretchr/testify/assert"
)

var imageController = NewImageController(configs.Client)

func seed_images() {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Document to insert
	images := []models.Image{
		// Define your test data here
		{
			CorId:             "1a",
			CreatorName:       "Bob",
			ImageName:         "image1",
			ImageTag:          "v1.0-Bob",
			ImageRegistryLink: "registry.com/bob",
		},
		{
			CorId:             "2b",
			CreatorName:       "Alice",
			ImageName:         "image2",
			ImageTag:          "v1.1-Alice",
			ImageRegistryLink: "registry.com/alice",
		},
	}

	// Conver the list of images to a slice of interface
	var documents []interface{}
	for _, image := range images {
		documents = append(documents, image)
	}

	// Insert the document
	_, err := configs.OpenCollection(configs.Client, "image_builder").InsertMany(ctx, documents)
	if err != nil {
		fmt.Printf("Error inserting document: %v\n", err)
		return
	}
}

func TestGetAllImages(t *testing.T) {
	// Seed the test data into the DB
	seed_images()

	// Create a new instance of the Gin router
	r := gin.Default()

	// Define the endpoint for the test
	r.GET("/image", imageController.GetAllImages)

	// Create a test request
	req, _ := http.NewRequest("GET", "/image", nil)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	expectedResponse := `[{"corId":"1a","creatorName":"Bob","imageName":"image1","imageTag":"v1.0-Bob","imageRegistryLink":"registry.com/bob"},{"corId":"2b","creatorName":"Alice","imageName":"image2","imageTag":"v1.1-Alice","imageRegistryLink":"registry.com/alice"}]`
	assert.Equal(t, expectedResponse, w.Body.String())
}

func TestGetImageByCorId(t *testing.T) {
    // Create a new instance of the Gin router
	r := gin.Default()

    // Define the endpoint for the test
	r.GET("/image/:corId", imageController.GetImageByCorId)

	// Create a test request
	req, _ := http.NewRequest("GET", "/image/2b", nil)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	expectedResponse := `{"corId":"2b","creatorName":"Alice","imageName":"image2","imageTag":"v1.1-Alice","imageRegistryLink":"registry.com/alice"}`
	assert.Equal(t, expectedResponse, w.Body.String())
}

func TestGetImageByCreatorName(t *testing.T) {
    // Create a new instance of the Gin router
	r := gin.Default()

    // Define the endpoint for the test
	r.GET("/image/name/:creatorName", imageController.GetImageByCreatorName)

	// Create a test request
	req, _ := http.NewRequest("GET", "/image/name/Bob", nil)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	expectedResponse := `[{"corId":"1a","creatorName":"Bob","imageName":"image1","imageTag":"v1.0-Bob","imageRegistryLink":"registry.com/bob"}]`
	assert.Equal(t, expectedResponse, w.Body.String())
}