// +build integration

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"platform_api/configs"
	"platform_api/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var attemptController = NewAttemptController(configs.Client)

func seed_attempts() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// document to insert
	attempts := []models.Attempt{
		{
			ChallengeName: "1a",
			CreatorName: "Winsen",
			Participant: "Phris",
			Token: "t1",
			ImageRegistryLink: "registry.com/winsen",
			SSHkey: "sshkey1",
			Result: 11,
			IpAddress: "127.0.0.1",
			Port: "11223",
		},
	}

	attemptBody := []models.AttemptBody {
		{
			ChallengeName: "2b",
			CreatorName: "Beep",
			Participant: "Bop",
			Token: "t2",
			ImageRegistryLink: "registry.com/beep",
		},
	}

	var documents []interface{}
	for _, attempt := range attempts {
		documents = append(documents, attempt)
	}

	for _, body := range attemptBody {
		documents = append(documents, body)
	}

	// Insert the document
	_, err := configs.OpenCollection(configs.Client, "attempt").InsertMany(ctx, documents)
	if err != nil {
		fmt.Printf("Error inserting document: %v\n", err)
		return
	}
}

func TestGetOneAttemptByToken(t *testing.T) {
	seed_attempts()

	r := gin.Default()

	r.GET("/attempt/:token", attemptController.GetOneAttemptByToken)

	req, _ := http.NewRequest("GET", "/attempt/t2", nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedResponse := `{"challengeName":"2b","creatorName":"Beep","participant":"Bop","token":"t2","imageRegistryLink":"registry.com/beep","sshkey":"","result":0,"ipaddress":"","port":""}`
	assert.Equal(t, expectedResponse, w.Body.String())
}

func TestGetOneAttemptByToken_NotFound(t *testing.T) {
	seed_attempts()
	
	r := gin.Default()
	
	r.GET("/attempt/:token", attemptController.GetOneAttemptByToken)

	req, _ := http.NewRequest("GET", "/attempt/xxx", nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	expectedResponse := `{"code":404,"message":"Failed to retrieve attempt","error":"attempt not found"}`

	assert.Equal(t, expectedResponse, w.Body.String())
}

func TestSubmitAttemptByToken(t *testing.T) {
	seed_attempts()
	
	r := gin.Default()
	
	r.POST("/attempt/submit", attemptController.SubmitAttemptByToken)

	bodyContent := []byte(`{"token": "t1", "result": 66}`)

	req, _ := http.NewRequest("POST", "/attempt/submit", bytes.NewBuffer(bodyContent))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedResponse := `{"submitted":true}`

	assert.Equal(t, expectedResponse, w.Body.String())
}

func TestSubmitAttemptByToken_NotFound(t *testing.T) {
	seed_attempts()
	
	r := gin.Default()
	
	r.POST("/attempt/submit", attemptController.SubmitAttemptByToken)

	bodyContent := []byte(`{"token": "xx", "result": 66}`)

	req, _ := http.NewRequest("POST", "/attempt/submit", bytes.NewBuffer(bodyContent))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	expectedResponse := `{"code":404,"message":"Failed to update an attempt","error":"attempt not found"}`

	assert.Equal(t, expectedResponse, w.Body.String())
}

func TestSubmitAttemptByToken_BadRequest(t *testing.T) {
	seed_attempts()
	
	r := gin.Default()
	
	r.POST("/attempt/submit", attemptController.SubmitAttemptByToken)

	bodyContent := []byte(`{"token": "xx"}`)

	req, _ := http.NewRequest("POST", "/attempt/submit", bytes.NewBuffer(bodyContent))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}