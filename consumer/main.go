package main

import (
	"go-rbmq/rabbitmq"
	"go-rbmq/models"
	"log"
	"encoding/json"
	"os"

)


const (
    DirectExchange  = "orders_direct"
    TopicExchange   = "orders_topic"
    FanoutExchange  = "orders_fanout"
)

func main() {
	// Connect to RabbitMQ
	rmq, err := rabbitmq.NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rmq.Close()
	// Declare queues
	_, err = rmq.DeclareQueue("orders")
	if err != nil {
		log.Fatal("Failed to declare orders queue:", err)
	}

	_, err = rmq.DeclareQueue("email_notifications")
	if err != nil {
		log.Fatal("Failed to declare email queue:", err)
	}

	_, err = rmq.DeclareQueue("sms_notifications")
	if err != nil {
		log.Fatal("Failed to declare SMS queue:", err)
	}
	orderConsumer(rmq)

	err = rmq.DeclareExchange(DirectExchange, "direct")
	if err != nil {
		log.Fatal("Failed to declare direct exchange:", err)
	}
	err = rmq.DeclareExchange(TopicExchange, "topic")
	if err != nil {
		log.Fatal("Failed to declare topic exchange:", err)
	}

	err = rmq.DeclareExchange(FanoutExchange, "fanout")
	if err != nil {
		log.Fatal("Failed to declare fanout exchange:", err)
	}
}


func orderConsumer(rmq *rabbitmq.RabbitMQ) {
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


