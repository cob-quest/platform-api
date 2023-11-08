package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	PORT            string
	API_LISTEN_HOST string

	AMQP_PORT     string
	AMQP_HOSTNAME string
	AMQP_USERNAME string
	AMQP_PASSWORD string

	MONGO_URI string
)

func InitEnv() {
	// loads environment variables
	err := godotenv.Load("secrets/.env")
	if err != nil {
		panic("Error loading env file")
	}

	// rest api
	PORT = getEnv("PORT", "8080")
	API_LISTEN_HOST = getEnv("API_LISTEN_HOST", "0.0.0.0")

	// amqp
	AMQP_PORT = getEnv("AMQP_PORT", "5672")
	AMQP_HOSTNAME = getEnv("AMQP_HOSTNAME", "rabbitmq.default.svc.cluster.local")
	AMQP_USERNAME = getEnv("AMQP_USERNAME", "rabbit")
	AMQP_PASSWORD = getEnv("AMQP_PASSWORD", "rabbit")

	MONGO_URI = GetMongoURI()

}

func GetMongoURI() string {
	err := godotenv.Load("secrets/.env")
	if err != nil {
		panic("Error loading env file")
	}
	MONGO_USER := os.Getenv("MONGODB_USERNAME")
	MONGO_PASS := os.Getenv("MONGODB_PASSWORD")
	MONGO_HOSTNAME := os.Getenv("MONGODB_HOSTNAME")
    MONGO_URI := os.Getenv("MONGO_URI")
	if MONGO_URI == "" {
		MONGO_URI = fmt.Sprintf("mongodb://%s:%s@%s:%s", MONGO_USER, MONGO_PASS, MONGO_HOSTNAME, "27017")
	}
	return MONGO_URI
}

// get env with default if the value is empty
// getEnv("ENV_VAR", "default")
func getEnv(s ...string) string {

	if len(s) <= 0 {

		// only one arg, don't provide defaults
		return ""

	} else if val := os.Getenv(s[0]); len(s) >= 2 && val != "" {

		// two args and the env var provides empty value
		return val

	} else {

		return s[1]

	}

}
