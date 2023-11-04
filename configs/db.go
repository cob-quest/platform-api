package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = ConnectDB()

func ConnectDB() (client *mongo.Client) {
	fmt.Printf("Attempting connection with %s\n", GetMongoURI())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI(GetMongoURI()))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Printf("Connecting to MongoDB ...")
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	fmt.Println("Success!")
	fmt.Println("Pinging server ...")

	// Veriify the connection to the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping database %v\n", err)
	}
	fmt.Println("Success")
	fmt.Println("Initalising indexes ...")

	// initialise indexes
	InitIndexes(client)
	fmt.Println("Successfully initialising indexes")
	return client
}

func InitIndexes(client *mongo.Client) {
    // Index for `image_builder` collection
	imageCollection := OpenCollection(client, "image_builder")
    imageCorIdIndexModel := mongo.IndexModel{
        Keys: bson.D{
			{Key: "corId", Value: 1},
        },
        Options: options.Index().SetUnique(true),
    }
	imageCompositeIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "creatorName", Value: 1},
			{Key: "imageName", Value: 1},
			{Key: "imageTag", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
    imageIndexModels := []mongo.IndexModel{imageCorIdIndexModel, imageCompositeIndexModel}
	imageIndexCreated, err := imageCollection.Indexes().CreateMany(context.Background(), imageIndexModels)
	if err != nil {
		log.Fatal(err)
	}

    // Index for `challenge` collection
	challengeCollection := OpenCollection(client, "challenge")
	challengeCorIdIndexModel := mongo.IndexModel{
		Keys:    bson.D{
            {Key: "corId", Value: 1},
        },
		Options: options.Index().SetUnique(true),
	}

    challengeIndexModels := []mongo.IndexModel{challengeCorIdIndexModel}
	challengeIndexCreated, err := challengeCollection.Indexes().CreateMany(context.Background(), challengeIndexModels)
	if err != nil {
		log.Fatal(err)
	}

    // Index for `process_engine` collection
	processCollection := OpenCollection(client, "process_engine")

	processUniqueIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "corId", Value: 1},
            {Key: "timestamp", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	processIndexCreated, err := processCollection.Indexes().CreateOne(context.Background(), processUniqueIndex)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created Image Index %s\n", imageIndexCreated)
	fmt.Printf("Created Challenge Index %s\n", challengeIndexCreated)
	fmt.Printf("Created Engine Index %s\n", processIndexCreated)
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cob").Collection(collectionName)
	return collection
}
