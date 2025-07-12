package controllers

import (
	"encoding/json"
	"go-rbmq/database"
	"go-rbmq/models"
	"go-rbmq/rabbitmq"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//query for orders
	query := "INSERT INTO orders (user_email, user_phone, amount, product_id) VALUES ($1, $2, $3, $4) RETURNING id"

	err := database.Db.QueryRow(query, order.UserEmail, order.UserPhone, order.Amount, order.ProductID).Scan(&order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}
	orderBytes, err := json.Marshal(order)
	if err != nil {
		log.Println("Error marshaling order:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process order"})
		return
	}
	err = rabbitmq.Rbmq.Publish("orders", orderBytes)
	if err != nil {
		log.Println("Error publishing order:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// Get all orders
func GetOrders(c *gin.Context) {
	query := "SELECT id, user_email, user_phone, amount, product_id FROM orders"
	rows, err := database.Db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserEmail, &order.UserPhone, &order.Amount, &order.ProductID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order"})
			return
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, orders)
}
