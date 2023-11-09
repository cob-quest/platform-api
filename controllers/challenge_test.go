package controllers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"platform_api/models"

	// "go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/stretchr/testify/assert"
)

// Seed test data

func seed_challenges() {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Document to insert
	images := []models.Challenge{
		// Define your test data here
		{
			CorID:             "1a",
			ChallengeName:     "ChallengeOne",
			CreatorName:       "Bob",
			ImageName:         "image1",
			ImageTag:          "v1.0-Bob",
			ImageRegistryLink: "registry.com/bob",
			Duration:          70,
			Participants:      []string{"gab@smu.com.sg"},
		},
		{
			CorID:             "2b",
			ChallengeName:     "ChallengeTwo",
			CreatorName:       "Alice",
			ImageName:         "image2",
			ImageTag:          "v1.1-Alice",
			ImageRegistryLink: "registry.com/alice",
			Duration:          60,
			Participants:      []string{"ben@smu.com.sg"},
		},
	}

	// Conver the list of images to a slice of interface
	var documents []interface{}
	for _, image := range images {
		documents = append(documents, image)
	}

	// Insert the document
	_, err := challengeCollection.InsertMany(ctx, documents)
	if err != nil {
		fmt.Printf("Error inserting document: %v\n", err)
		return
	}
}

func TestGetAllChallenges(t *testing.T) {
	// Seed the test data into the DB
	seed_challenges()

	// Create a new instance of the Gin router
	r := gin.Default()

	// Create an image controller
	challengeController := new(ChallengeController)

	// Define the endpoint for the test
	r.GET("/challenge", challengeController.GetAllChallenges)

	// Create a test request
	req, _ := http.NewRequest("GET", "/challenge", nil)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	expectedResponse := `[{"corId":"1a","challengeName":"ChallengeOne","creatorName":"Bob","imageName":"image1","imageTag":"v1.0-Bob","imageRegistryLink":"registry.com/bob","duration":70,"participants":["gab@smu.com.sg"]},{"corId":"2b","challengeName":"ChallengeTwo","creatorName":"Alice","imageName":"image2","imageTag":"v1.1-Alice","imageRegistryLink":"registry.com/alice","duration":60,"participants":["ben@smu.com.sg"]}]`
	assert.Equal(t, expectedResponse, w.Body.String())
}

func TestGetChallengeByCorId(t *testing.T) {
    // Create a new instance of the Gin router
	r := gin.Default()

	// Create an image controller
	challengeController := new(ChallengeController)

    // Define the endpoint for the test
	r.GET("/challenge/:corId", challengeController.GetChallengeByCorID)

	// Create a test request
	req, _ := http.NewRequest("GET", "/challenge/2b", nil)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	expectedResponse := `{"corId":"2b","challengeName":"ChallengeTwo","creatorName":"Alice","imageName":"image2","imageTag":"v1.1-Alice","imageRegistryLink":"registry.com/alice","duration":60,"participants":["ben@smu.com.sg"]}`
	assert.Equal(t, expectedResponse, w.Body.String())
}

func TestGetChallengeByCreatorName(t *testing.T) {
    // Create a new instance of the Gin router
	r := gin.Default()

	// Create an image controller
	challengeController := new(ChallengeController)

    // Define the endpoint for the test
	r.GET("/challenge/:creatorName", challengeController.GetChallengeByCreatorName)

	// Create a test request
	req, _ := http.NewRequest("GET", "/challenge/Alice", nil)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	expectedResponse := `[{"corId":"2b","challengeName":"ChallengeTwo","creatorName":"Alice","imageName":"image2","imageTag":"v1.1-Alice","imageRegistryLink":"registry.com/alice","duration":60,"participants":["ben@smu.com.sg"]}]`
	assert.Equal(t, expectedResponse, w.Body.String())
}