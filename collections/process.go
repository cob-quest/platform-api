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

type ProcessCollection struct {
	Collection *mongo.Collection
}

func NewProcessCollection(client *mongo.Client) *ProcessCollection {
	return &ProcessCollection{Collection: configs.OpenCollection(client, "process_engine")}
}

func (t ProcessCollection) GetAllProcesses() (*[]models.Process, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find()
	cursor, err := t.Collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var process []models.Process
	err = cursor.All(ctx, &process)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &process, http.StatusOK, nil
}

func (t ProcessCollection) GetAllProcessByCorID(corId string) (*[]models.Process, int, error) {
	if corId == "" {
		return nil, http.StatusBadRequest, errors.New("corId cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "corId", Value: corId}}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}) // -1 for descending order
	cursor, err := t.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	defer cursor.Close(ctx)

	var process []models.Process
	err = cursor.All(ctx, &process)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(process) == 0 {
		return nil, http.StatusNotFound, errors.New("process is not found by the corId")
	}

	return &process, http.StatusOK, nil
}

func (t ProcessCollection) GetProcessByCreatorName(creatorName string) (*[]models.Process, int , error) {
	if creatorName == "" {
		return nil, http.StatusBadRequest, errors.New("creatorName cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "creatorName", Value: creatorName}}
	cursor, err := t.Collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}

	defer cursor.Close(ctx)

	var process []models.Process
	err = cursor.All(ctx, &process)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(process) == 0 {
		return nil, http.StatusNotFound, errors.New("no process found with creatorName")
	}

	return &process, http.StatusOK, nil
}

func (t ProcessCollection) GetLatestStatusByCorId(corId string) (*models.Process, int, error){
	if corId == "" {
		return nil, http.StatusBadRequest, errors.New("corId cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var process models.Process
	filter := bson.D{
		{Key: "corId", Value: corId},
	}

	options := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	err := t.Collection.FindOne(ctx, filter, options).Decode(&process)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, http.StatusNotFound, errors.New("no process found with corId")
		} else {
			return nil, http.StatusInternalServerError, err
		}
	}

	return &process, http.StatusOK, nil
}