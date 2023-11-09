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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChallengeCollection struct {
	Collection *mongo.Collection
}

func NewChallengeCollection(client *mongo.Client) *ChallengeCollection {
	return &ChallengeCollection{Collection: configs.OpenCollection(client, "challenge")}
}

func (t ChallengeCollection) GetAllChallenges() (*[]models.Challenge, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find()
	cursor, err := t.Collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var challenges []models.Challenge
	err = cursor.All(ctx, &challenges)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &challenges, http.StatusOK, nil
}

func (t ChallengeCollection) GetChallengeByCorID(corId string) (*models.Challenge, int, error) {
	if corId == "" {
		return nil, http.StatusBadRequest, errors.New("corId cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var challenge models.Challenge

	filter := bson.D{{Key: "corId", Value: corId}}
	err := t.Collection.FindOne(ctx, filter).Decode(&challenge)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, http.StatusNotFound, errors.New("no challenge found with given challenge corId")
		} else {
			return nil, http.StatusInternalServerError, err
		}
	}

	return &challenge, http.StatusOK, nil
}

func (t ChallengeCollection) GetChallengeByCreatorName(creatorName string) (*[]models.Challenge, int, error) {
	if creatorName == "" {
		return nil, http.StatusBadRequest, errors.New("creator name cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "creatorName", Value: creatorName}}
	cursor, err := t.Collection.Find(ctx, filter)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	defer cursor.Close(ctx)

	var challenges []models.Challenge

	err = cursor.All(ctx, &challenges)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(challenges) == 0 {
		return nil, http.StatusNotFound, errors.New("no challenges with creator name found")
	}

	return &challenges, http.StatusOK, nil
}

func (t ChallengeCollection) CheckChallengeByChallengeAndCreatorName(challengeName string, creatorName string) (int, error) {
	if creatorName == "" || challengeName == "" {
		return http.StatusBadRequest, errors.New("creator and challenge name cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var challenge models.Challenge
	filter := bson.D{{
		Key:   "challengeName",
		Value: challengeName,
	}, {
		Key:   "creatorName",
		Value: creatorName,
	}}

	err := t.Collection.FindOne(ctx, filter).Decode(&challenge)
	if err == nil {
		return http.StatusBadRequest, errors.New("challenge name already exists")
	}

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}