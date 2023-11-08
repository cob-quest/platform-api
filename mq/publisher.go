package mq

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Pub (Publishes) a specified message to the AMQP exchange
func Pub(ex string, key string, body []byte) error {

	// open channel
	c, err := AMQP_CONN.Channel()
	for err != nil {
		log.Printf("Failed to connect: %s, retrying in 5 seconds...", err)
		c, err = ReestablishConnection()
		time.Sleep(5 * time.Second)
	}
	defer c.Close()

	// context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// publish message
	err = c.PublishWithContext(
		ctx,
		ex,
		key,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Sending message: %s -> %s", body, "default")
	return nil
}

