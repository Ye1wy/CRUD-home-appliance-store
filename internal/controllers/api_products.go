package controllers

import (
	"github.com/gin-gonic/gin"
)

type ProductsAPI struct {
}

// Post /api/v1/api/v1/products
// Add product
func (api *ProductsAPI) AddProduct(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Patch /api/v1/api/v1/products/:id/decrease
// Decrease a parameter by a given value
func (api *ProductsAPI) DecreaseParametr(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Delete /api/v1/api/v1/products/:id
// Delete product by identificator
func (api *ProductsAPI) DeleteProductById(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Get /api/v1/api/v1/products
// Retrieve all product
func (api *ProductsAPI) GetAllProduct(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Get /api/v1/api/v1/products/:id
// Get product by id
func (api *ProductsAPI) SearchProductById(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}
