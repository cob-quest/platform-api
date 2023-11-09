package collections

import (
	"context"
	"errors"
	"net/http"
	"platform_api/configs"
	"platform_api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AttemptCollection struct {
	Collection *mongo.Collection
}

func NewAttemptCollection(client *mongo.Client) *AttemptCollection {
	return &AttemptCollection{Collection: configs.OpenCollection(client, "attempt")}
}

func (t AttemptCollection) GetAllAttempts() (*[]models.Attempt, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := t.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var attempts []models.Attempt
	err = cursor.All(ctx, &attempts)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &attempts, http.StatusOK, nil
}

func (t AttemptCollection) GetOneAttemptByToken(token string) (*models.Attempt, int, error) {
	if token == "" {
		return nil, http.StatusBadRequest, errors.New("invalid token parameter")
	}

	// ctx for 10s
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find the challenge with the specified ID
	var attemptSingle models.Attempt
	filter := bson.D{
		{Key: "token", Value: token},
		// {Key: "email", Value: req.Email},
	}
	err := t.Collection.FindOne(ctx, filter).Decode(&attemptSingle)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, http.StatusNotFound, errors.New("attempt not found")
		} else {
			return nil, http.StatusInternalServerError, err
		}
	}

	return &attemptSingle, http.StatusOK, nil
}

func (t AttemptCollection) UpdateAnAttempt(attemptSubmitBody *models.AttemptSubmitBody, token string) (int, error) {
	if token == "" {
		return http.StatusBadRequest, errors.New("invalid token parameter")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "token", Value: token}}

	// Create an update document to update the value of the object.
	var attempt models.Attempt
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "result", Value: attemptSubmitBody.Result}}}}
	err := t.Collection.FindOneAndUpdate(ctx, filter, update).Decode(&attempt)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return http.StatusNotFound, errors.New("attempt not found")
		} else {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}