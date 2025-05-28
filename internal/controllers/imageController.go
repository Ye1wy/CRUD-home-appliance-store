package controllers

import (
	crud_errors "CRUD-HOME-APPLIANCE-STORE/internal/errors"
	"CRUD-HOME-APPLIANCE-STORE/internal/mapper"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type imageService interface {
	Create(ctx context.Context, image *domain.Image) error
	GetAll(ctx context.Context, limit, offset int) ([]domain.Image, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Image, error)
	Update(ctx context.Context, image *domain.Image) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ImageController struct {
	*BaseController
	service imageService
}

func NewImageController(service imageService, logger *logger.Logger) *ImageController {
	base := NewBaseContorller(logger)
	return &ImageController{
		base,
		service,
	}
}

func (ctrl *ImageController) Create(c *gin.Context) {
	op := "controllers.imageController.Create"
	imageRaw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		ctrl.logger.Error("Failed to read image bytes", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid image body"})
		return
	}

	image := domain.Image{Data: imageRaw}

	if err := ctrl.service.Create(c.Request.Context(), &image); err != nil {
		if errors.Is(err, crud_errors.ErrImageCorruption) {
			ctrl.logger.Error("Invalid image data", logger.Err(err), "op", op)
			ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "invalid image"})
			return
		}

		ctrl.logger.Error("Failed to add image", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"massage": "server is busy"})
		return
	}

	ctrl.logger.Debug("Image added successfully", "op", op)
	c.Status(http.StatusCreated)
}

func (ctrl *ImageController) GetAll(c *gin.Context) {
	op := "controllers.imageController.GetAll"
	limit, err := strconv.Atoi(c.DefaultQuery("limit", defaultLimit))
	if err != nil {
		ctrl.logger.Error("Invalid limit parameter", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("limit", defaultLimit))
	if err != nil {
		ctrl.logger.Error("Invalid offset parameter", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{})
		return
	}

	images, err := ctrl.service.GetAll(c.Request.Context(), limit, offset)
	if errors.Is(err, crud_errors.ErrNotFound) {
		ctrl.logger.Warn("Images not found", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusNoContent, gin.H{"massage": "no have images"})
		return
	}

	if err != nil {
		ctrl.logger.Error("Failed to retrive images", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Failed to retrieve client"})
		return
	}

	imagesDTOs := make([]dto.Image, len(images), cap(images))

	for i, image := range images {
		dto := mapper.ImageToDTO(image)
		imagesDTOs[i] = dto
	}

	ctrl.logger.Info("Retrieved all images", "limit", limit, "offset", offset, "op", op)
	ctrl.responce(c, http.StatusOK, gin.H{})
}

func (ctrl *ImageController) GetById(c *gin.Context) {
	op := "controllers.imageController.GetById"
	rawId := c.Query("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Error("Invaldid taken id", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload"})
		return
	}

	image, err := ctrl.service.GetById(c.Request.Context(), id)
	if err != nil {
		ctrl.logger.Error("Failed to get data from database", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"massage": "server is busy"})
		return
	}

	dto := mapper.ImageToDTO(*image)

	ctrl.logger.Info("Image data is retrived")
	ctrl.responce(c, http.StatusOK, dto)
}

func (ctrl *ImageController) Update(c *gin.Context) {
	op := "controllers.imageController.Delete"
	rawdId := c.Param("id")
	id, err := uuid.Parse(rawdId)
	if err != nil {
		ctrl.logger.Warn("Invalid id", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload"})
		return
	}

	imageRaw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		ctrl.logger.Error("Failed to read image bytes", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid image body"})
		return
	}

	image := domain.Image{
		Id:   id,
		Data: imageRaw,
	}

	if err := ctrl.service.Update(c.Request.Context(), &image); err != nil {
		if errors.Is(err, crud_errors.ErrImageCorruption) {
			ctrl.logger.Error("Invalid image data", logger.Err(err), "op", op)
			ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "invalid image"})
			return
		}

		ctrl.logger.Error("Failed to update image", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"massage": "server is busy"})
		return
	}

	ctrl.logger.Debug("Image is updated", "op", op)
	c.Status(http.StatusOK)
}

func (ctrl *ImageController) Delete(c *gin.Context) {
	op := "controllers.imageController.Delete"
	rawId := c.Param("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Warn("Invalid id", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload"})
		return
	}

	if err := ctrl.service.Delete(c.Request.Context(), id); err != nil {
		ctrl.logger.Error("Failed to delete image from database", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"massage": "server is busy"})
		return
	}

	ctrl.logger.Debug("Image deleted", "op", op)
	c.Status(http.StatusNoContent)
}
