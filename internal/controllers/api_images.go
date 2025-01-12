package controllers

import (
	"github.com/gin-gonic/gin"
)

type ImagesAPI struct {
}

// Post /api/v1/api/v1/images
// Add image
func (api *ImagesAPI) AddImage(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Patch /api/v1/api/v1/images/:id/changeImage
// Change a image
func (api *ImagesAPI) ChangeImage(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Delete /api/v1/api/v1/images/:id
// Delete image by identificator
func (api *ImagesAPI) DeleteImageById(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Get /api/v1/api/v1/images/:id
// Get image by id
func (api *ImagesAPI) SearchImageById(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Get /api/v1/api/v1/images/products/:id
// Getting an image of a specific product
func (api *ImagesAPI) SearchProductImage(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}
