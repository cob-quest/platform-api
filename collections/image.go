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

type ImageCollection struct {
	Collection *mongo.Collection
}

func NewImageCollection(client *mongo.Client) *ImageCollection {
	return &ImageCollection{Collection: configs.OpenCollection(client, "image_builder")}
}

func (t ImageCollection) GetAllImages() (*[]models.Image, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find()
	cursor, err := t.Collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	defer cursor.Close(ctx)

	var images []models.Image
	err = cursor.All(ctx, &images)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &images, http.StatusOK, nil
}

func (t ImageCollection) GetImageByCorId(corId string) (*models.Image, int, error) {
	if corId == "" {
		return nil, http.StatusBadRequest, errors.New("corId cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var image models.Image

	filter := bson.D{{Key: "corId", Value: corId}}
	err := t.Collection.FindOne(ctx, filter).Decode(&image)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, http.StatusNotFound, errors.New("no image found with given challenge corId")
		} else {
			return nil, http.StatusInternalServerError, err
		}
	}

	return &image, http.StatusOK, nil
}

func (t ImageCollection) GetImageByCreatorName(creatorName string) (*[]models.Image, int, error) {
	if creatorName == "" {
		return nil, http.StatusBadRequest, errors.New("creatorName cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "creatorName", Value: creatorName}}
	cursor, err := t.Collection.Find(ctx, filter)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	defer cursor.Close(ctx)

	var images []models.Image
	err = cursor.All(ctx, &images)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(images) == 0 {
		return nil, http.StatusNotFound, errors.New("no images with creatorName found")
	}

	return &images, http.StatusOK, nil
}

func (t ImageCollection) CheckImageByImageAndCreatorName(imageName string, creatorName string) (int, error) {
	if creatorName == "" || imageName == "" {
		return http.StatusBadRequest, errors.New("creator and image name cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var challenge models.Challenge
	filter := bson.D{{
		Key:   "imageName",
		Value: imageName,
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

// ------- FOR CHALLENGE CONTROLLER ---------
func (t ImageCollection) CheckImageExists(imageName string, imageTag string, creatorName string) (int, error) {
	if imageName == "" || imageTag == "" || creatorName == "" {
		return http.StatusBadRequest, errors.New("image name, tag, and creatorName cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var image models.Image
	filter := bson.D{
		{Key: "imageName", Value: imageName},
		{Key: "creatorName", Value: creatorName},
		{Key: "imageTag", Value: imageTag},
	}
	err := t.Collection.FindOne(ctx, filter).Decode(&image)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return http.StatusNotFound, errors.New("image is not found given image name, tag, and creatorName")
		} else {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}
