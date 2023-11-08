package controllers

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "platform_api/models"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockImageCollection struct {
    mock.Mock
}

func (m *MockImageCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
    args := m.Called(ctx, filter, opts)
    return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockImageCollection) All(ctx context.Context, result interface{}) error {
    args := m.Called(ctx, result)
    return args.Error(0)
}

func TestGetAllImages(t *testing.T) {
    // Initialize the mock database collection
    mockCollection := new(MockImageCollection)

    // Create a new instance of the Gin router
    r := gin.Default()

    // Create an instance of ImageController with the mock collection
    imageController := ImageController{}

    // Define test data and expected results
    mockCursor := new(mongo.Cursor)
    expectedImages := []models.Image{
        // Define your test data here
        {
            CorId: "1a",
            CreatorName: "Bob",
            ImageName: "image1",
            ImageTag: "v1.0-Bob",
            ImageRegistryLink: "registry.com/bob",
        },
        {
            CorId: "2b",
            CreatorName: "Alice",
            ImageName: "image2",
            ImageTag: "v1.1-Alice",
            ImageRegistryLink: "registry.com/alice",
        },
    }

    // Stub the Find method to return the mock cursor
    mockCollection.On("Find", mock.Anything, bson.M{}, mock.AnythingOfType("*options.FindOptions")).
        Return(mockCursor, nil)

    // Stub the All method to populate the expectedImages slice
    mockCollection.On("All", mock.Anything, &expectedImages).
        Return(nil)

    // Create a test request
    req, _ := http.NewRequest("GET", "/image", nil)

    // Create a test response recorder
    w := httptest.NewRecorder()

    // Define the endpoint for the test
    r.GET("/image", imageController.GetAllImages)

    // Perform the request
    r.ServeHTTP(w, req)

    // Check the response status code
    assert.Equal(t, http.StatusOK, w.Code)

    // Check the response body
    expectedResponse := `{"id":1,"name":"Image1"},{"id":2,"name":"Image2"}` // Will edit the expected ltr
    assert.Equal(t, expectedResponse, w.Body.String())

    // Assert that the Find and All methods were called as expected
    mockCollection.AssertExpectations(t)
}
