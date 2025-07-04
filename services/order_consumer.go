package main

import (
	"encoding/json"
	"log"

	"go-rbmq/models"
	"go-rbmq/rabbitmq"
)

func main() {
	// Connect to RabbitMQ
	rmq, err := rabbitmq.NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rmq.Close()

	// Consume from orders queue
	msgs, err := rmq.Consume("orders")
	if err != nil {
		log.Fatal("Failed to consume from queue:", err)
	}

	log.Println("Order Processing Service started. Waiting for orders...")

	for msg := range msgs {
		var order models.Order
		err := json.Unmarshal(msg.Body, &order)
		if err != nil {
			log.Println("Error decoding order:", err)
			msg.Nack(false, true) // Requeue message
			continue
		}

		// Process the order (in a real system, this would do more)
		log.Printf("Processing order %s for %s", order.ID, order.UserEmail)

		// Send notifications
		sendNotifications(rmq, order)

		msg.Ack(false) // Acknowledge message
	}
}

func sendNotifications(rmq *rabbitmq.RabbitMQ, order models.Order) {
	orderBytes, err := json.Marshal(order)
	if err != nil {
		log.Println("Error marshaling order for notifications:", err)
		return
	}

	// Send email notification
	err = rmq.Publish("email_notifications", orderBytes)
	if err != nil {
		log.Println("Error publishing email notification:", err)
	}

}
