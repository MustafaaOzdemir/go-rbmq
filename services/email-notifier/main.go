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

	// Consume from email notifications queue
	msgs, err := rmq.Consume("email_notifications")
	if err != nil {
		log.Fatal("Failed to consume from queue:", err)
	}

	log.Println("Email Notifier Service started. Waiting for messages...")

	for msg := range msgs {
		var order models.Order
		err := json.Unmarshal(msg.Body, &order)
		if err != nil {
			log.Println("Error decoding order:", err)
			msg.Nack(false, true) // Requeue message
			continue
		}

		// Process email notification
		log.Printf("Sending email to %s for order %s (Amount: $%.2f)\n",
			order.UserEmail, order.ID, order.Amount)

		// Simulate email sending
		// In real world, integrate with SendGrid/Mailgun/etc.

		msg.Ack(false) // Acknowledge message
	}
}
