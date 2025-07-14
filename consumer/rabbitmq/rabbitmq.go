package rabbitmq

import (
	"time"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

var Rbmq RabbitMQ

func NewRabbitMQ() (*RabbitMQ, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}

	user := os.Getenv("RBMQUSER")

	password := os.Getenv("RBMQPASSWORD")


	// set up postgres sql to open it.
	url := fmt.Sprintf("amqp://%s:%s@localhost:5672/",
		user, password)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	Rbmq = RabbitMQ{Conn: conn, Channel: channel}

	return &RabbitMQ{Conn: conn, Channel: channel}, nil
}

func (r *RabbitMQ) Close() {
	r.Channel.Close()
	r.Conn.Close()
}
func (r *RabbitMQ) DeclareExchange(name, kind string) error {
	return r.Channel.ExchangeDeclare(
		name,
		kind, // type
		true,
		false,
		false,
		false,
		nil,
	)
}
func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.Channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func (r *RabbitMQ) Publish(queue string, body []byte) error {
	return r.Channel.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
}

func (r *RabbitMQ) Consume(queue string) (<-chan amqp.Delivery, error) {
	return r.Channel.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

