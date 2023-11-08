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

func seed() {
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
