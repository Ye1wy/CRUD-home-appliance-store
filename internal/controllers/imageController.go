package controllers

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/services"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ImagesController struct {
	service services.ImagesService
	logger  *slog.Logger
}

func NewImagesController(imagesService services.ImagesService, log *slog.Logger) *ImagesController {
	return &ImagesController{
		service: imagesService,
		logger:  log,
	}
}

// Post /api/v1/images
// Add image
func (ctrl *ImagesController) AddImage(c *gin.Context) {
	op := "controllers.image.addImage"

	var image model.Image
	if err := c.ShouldBindJSON(&image); err != nil {
		ctrl.logger.Error("Failed to bind JSON for AddImage", logger.Err(err), "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ctrl.service.AddImage(c.Request.Context(), &image); err != nil {
		ctrl.logger.Error("Failed to add image: ", logger.Err(err), "op", op)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add image"})
		return
	}

	ctrl.logger.Info("Image added successfully", "imageID", image.Id, "op", op)
	c.JSON(http.StatusCreated, gin.H{"message": "Image added successfully"})
}

// Get /api/v1/images/products/:id
// Getting an image of a specific product
func (ctrl *ImagesController) SearchProductImage(c *gin.Context) {
	op := "controllers.image.searchProductImage"

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ctrl.logger.Error("Invalid product ID parameter for SearchProductImage", logger.Err(err), "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	image, err := ctrl.service.GetProductImages(c.Request.Context(), id)
	if image == nil {
		ctrl.logger.Warn("Product not found", logger.Err(err), "op", op)
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err != nil {
		ctrl.logger.Error("Error while searching image ", logger.Err(err), "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to found image"})
		return
	}

	ctrl.logger.Info("Image retrieved successefully", "id", id, "op", op)
	c.JSON(http.StatusOK, image)

}

// Get /api/v1/images/:id
// Get image by id
func (ctrl *ImagesController) SearchImageById(c *gin.Context) {
	op := "controllers.image.getImageById"

	id := c.Query("id")

	image, err := ctrl.service.GetImageById(c.Request.Context(), id)
	if image == nil && err == nil {
		ctrl.logger.Warn("Image not found", logger.Err(err), "op", op)
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	if err != nil {
		ctrl.logger.Error("Error while searching image ", logger.Err(err), "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to found image"})
		return
	}

	ctrl.logger.Info("Image retrieved successfully", "id", id, "op", op)
	c.JSON(http.StatusOK, image)
}

// Patch /api/v1/images/:id/changeImage
// Change a image
func (ctrl *ImagesController) ChangeImage(c *gin.Context) {
	op := "controllers.image.changeImage"

	id := c.Param("id")

	var updatedFields model.ChangeImageRequest

	if err := c.ShouldBindJSON(&updatedFields); err != nil {
		ctrl.logger.Error("Failed to bind JSON for ChangeImage", logger.Err(err), "op", op)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ctrl.service.ChangeImage(c.Request.Context(), id, updatedFields.Image); err != nil {
		ctrl.logger.Error("Failed to update address ID", "imageID", id, logger.Err(err), "op", op)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address ID"})
		return
	}

	ctrl.logger.Info("Address ID updated successfully", "imageID", id, "op", op)
	c.JSON(http.StatusOK, gin.H{"status": "Address ID updated"})
}

// Delete /api/v1/images/:id
// Delete image by identificator
func (ctrl *ImagesController) DeleteImageById(c *gin.Context) {
	op := "controllers.image.deleteImageById"

	id := c.Param("id")

	if err := ctrl.service.DeleteImageById(c.Request.Context(), id); err != nil {
		ctrl.logger.Error("Failed to delete image", "imageID", id, logger.Err(err), "op", op)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
		return
	}

	ctrl.logger.Info("Image deleted successfully", "imageID", id, "op", op)
	c.Status(http.StatusNoContent)
}
