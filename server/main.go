package main

import (
	"go-rbmq/controllers"
	"go-rbmq/database"
	"go-rbmq/rabbitmq"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
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
	route.GET("/products", controllers.GetProducts)
	route.POST("/products", controllers.CreateProduct)
	route.GET("/products/:id", controllers.GetProductByID)
	route.PUT("/products/:id", controllers.UpdateProduct)
	route.DELETE("/products/:id", controllers.DeleteProduct)

	route.GET("/orders", controllers.GetOrders)
	route.POST("/orders", controllers.CreateOrder)

	database.ConnectDatabase()

	err = route.Run(":4200")
	if err != nil {
		panic(err)
	}

}
