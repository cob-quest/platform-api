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
	fmt.Println(EnvMongoUri())
	fmt.Printf("Attempting connection with %s\n", EnvMongoUri())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoUri()))
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

func InitIndexes(client * mongo.Client) {
	challengeCollection := OpenCollection(client, "challenge")
	challengeIndexModel := mongo.IndexModel{
        Keys: bson.D{{Key:"image_name", Value: 1}, {Key:"image_ver", Value: -1}},
		Options: options.Index().SetUnique(true),
	}
	challengeIndexCreated, err := challengeCollection.Indexes().CreateOne(context.Background(), challengeIndexModel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created Challenge Index %s\n", challengeIndexCreated)
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cs302").Collection(collectionName)
	return collection
}
