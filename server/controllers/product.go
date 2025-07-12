package controllers

import (
	"go-rbmq/database"
	"go-rbmq/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create a new product
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO products (name, description, price) VALUES ($1, $2, $3) RETURNING id"
	err := database.Db.QueryRow(query, product.Name, product.Description, product.Price).Scan(&product.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// Get all products
func GetProducts(c *gin.Context) {
	query := "SELECT id, name, description, price FROM products"
	rows, err := database.Db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan item"})
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

// Get a product by ID
func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	query := "SELECT id, name, description, price FROM products WHERE id = $1"
	var product models.Product
	err := database.Db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Update a product by ID
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4"
	_, err := database.Db.Exec(query, product.Name, product.Description, product.Price, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
}

// Delete a product by ID
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	query := "DELETE FROM products WHERE id = $1"
	_, err := database.Db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
