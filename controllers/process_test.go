// +build integration
package controllers

import (
	"context"
	"encoding/json"
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

var processController = NewProcessController(configs.Client)


func seed_processes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var processes []models.Process
	
	var dataJSON = `[{"corId": "id1","creatorName": "Ama","challengeName": "challenge1","imageName": "image1","imageTag": "latest","participants": ["ama@a.com"]},{"corId": "id2","creatorName": "Bob","challengeName": "challenge2","imageName": "image2","imageTag": "latest","participants": ["bob@b.com"]}]`
	err := json.Unmarshal([]byte(dataJSON), &processes)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(processes)

	var documents []interface{}
	for _, process := range processes {
		documents = append(documents, process)
	}

	// Insert the document
	_, err = configs.OpenCollection(configs.Client, "process_engine").InsertMany(ctx, documents)
	if err != nil {
		fmt.Printf("Error inserting document: %v\n", err)
		return
	}
}

func TestGetAllProcesses(t *testing.T) {
	seed_processes()

	r := gin.Default()

	r.GET("/process", processController.GetAllProcesses)

	req, _ := http.NewRequest("GET", "/process", nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	expectedResponse := `[{"timestamp":null,"corId":"id1","event":null,"eventStatus":null,"creatorName":"Ama","challengeName":"challenge1","imageName":"image1","imageTag":"latest","participants":["ama@a.com"]},{"timestamp":null,"corId":"id2","event":null,"eventStatus":null,"creatorName":"Bob","challengeName":"challenge2","imageName":"image2","imageTag":"latest","participants":["bob@b.com"]}]`

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

func TestGetProcessByCorId(t *testing.T) {
	seed_processes()

	r := gin.Default()

	r.GET("/process/:corId", processController.GetProcessByCorID)

	req, _ := http.NewRequest("GET", "/process/id2", nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	expectedResponse := `[{"timestamp":null,"corId":"id2","event":null,"eventStatus":null,"creatorName":"Bob","challengeName":"challenge2","imageName":"image2","imageTag":"latest","participants":["bob@b.com"]}]`
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

func TestGetProcessByCorId_NotFound(t *testing.T) {
	seed_processes()

	r := gin.Default()

	r.GET("/process/:corId", processController.GetProcessByCorID)

	req, _ := http.NewRequest("GET", "/process/xyxy", nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
}