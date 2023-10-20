package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoUri() string {
	err := godotenv.Load("./secrets/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	MONGO_USER := os.Getenv("MONGO_USER")
	MONGO_PASS := os.Getenv("MONGO_PASS")
	MONGO_HOSTNAME := os.Getenv("MONGO_HOSTNAME")
	return fmt.Sprintf("mongodb://%s,%s@%s:%s", MONGO_USER, MONGO_PASS, MONGO_HOSTNAME, "27017")
}
