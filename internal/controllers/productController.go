package controllers

import (
	"github.com/gin-gonic/gin"
)

type ProductsController struct {
}

// Post /api/v1/products
// Add product
func (ctrl *ProductsController) AddProduct(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Patch /api/v1/products/:id/decrease
// Decrease a parameter by a given value
func (ctrl *ProductsController) DecreaseParametr(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Delete /api/v1/products/:id
// Delete product by identificator
func (ctrl *ProductsController) DeleteProductById(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Get /api/v1/products
// Retrieve all product
func (ctrl *ProductsController) GetAllProduct(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Get /api/v1/products/:id
// Get product by id
func (ctrl *ProductsController) SearchProductById(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}
