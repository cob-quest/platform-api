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
	imageCollection := OpenCollection(client, "image_builder")
	imageIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "corId", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	imageIndexCreated, err := imageCollection.Indexes().CreateOne(context.Background(), imageIndexModel)
	if err != nil {
		log.Fatal(err)
	}

	challengeCollection := OpenCollection(client, "challenge_builder")
	challengeIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "corId", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	challengeIndexCreated, err := challengeCollection.Indexes().CreateOne(context.Background(), challengeIndexModel)
	if err != nil {
		log.Fatal(err)
	}

	processCollection := OpenCollection(client, "process_engine")

	processUniqueIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "corId", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	processEventTextModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "event", Value: "text"},
		},
	}

	processIndexes := []mongo.IndexModel{processUniqueIndex, processEventTextModel}

	processIndexCreated, err := processCollection.Indexes().CreateMany(context.Background(), processIndexes)
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
