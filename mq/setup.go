package mq

import (
	"fmt"
	"log"
	"platform_api/configs"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	AMQP_CONN *amqp.Connection
)

const (
	ROUTE_IMAGE_BUILD      = "platform.fromService.imageCreate"
	ROUTE_CHALLENGE_CREATE = "platform.fromService.challengeCreate"
	ROUTE_CHALLENGE_START  = "platform.fromService.challengeStart"
	QUEUE_PLATFORM_FROM     = "queue.platform.fromService"
	EXCHANGE_TOPIC_ROUTER = "topic.router"
	EXCHANGE_DEFAULT       = "/"
)

// Init rabbitmq connection
func Init() {

	// connect to mq
	connStr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		configs.AMQP_USERNAME,
		configs.AMQP_PASSWORD,
		configs.AMQP_HOSTNAME,
		configs.AMQP_PORT,
	)

	// setup rabbitmq connection
	var err error
	AMQP_CONN, err = amqp.Dial(connStr)
	log.Printf(
		"Connecting to host with %s",
		connStr,
	)
	if err != nil {
		log.Panic("Failed to connect to RabbitMQ", err)
	}
	log.Printf(
		"Successfully connected to RabbitMQ: %s:%s",
		configs.AMQP_HOSTNAME,
		configs.AMQP_PORT,
	)

	// conn to channel
	channel, err := AMQP_CONN.Channel()
	if err != nil {
		log.Panic("Failed to create channel for RabbitMQ", err)
	}

	// declare all queues
	QueueDeclare(channel)
	ExchangeDeclare(channel)
	channel.Close()

}

// Declares an exchange to Pub/Sub to
func ExchangeDeclare(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		EXCHANGE_TOPIC_ROUTER, // name
		"topic",                // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
}

// Declares a queue to Pub/Sub to
func QueueDeclare(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		QUEUE_PLATFORM_FROM, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
}
