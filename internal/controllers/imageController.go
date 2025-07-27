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

// func getImageBuffer(file multipart.File) ([]byte, error) {
// 	buf := bytes.NewBuffer(nil)
// 	if _, err := io.Copy(buf, file); err != nil {
// 		return nil, err
// 	}

// 	return buf.Bytes(), nil
// }

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

// CreateImage godoc
//
//	@Summary		Create image
//	@Description	Image created from JSON or XML, for create endpoint required: image
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			image	body	domain.Image	true	"Image data"
//	@Success		201
//	@Failure		400	{object}	domain.Error
//	@Failure		404	{object}	domain.Error
//	@Failure		500	{object}	domain.Error
//	@Router			/api/v1/images [post]
func (ctrl *ImageController) Create(c *gin.Context) {
	op := "controllers.imageController.Create"

	ctrl.logger.Debug("Content-Type contains", "data", c.ContentType())

	if c.ContentType() != contentTypeOctetStream {
		ctrl.logger.Warn("Invalid content-type", "got", c.ContentType(), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"message": "Expected multipart/form-data"})
		return
	}

	title := c.Request.Header.Get(headerXImageTitle)

	if title == "" {
		ctrl.logger.Warn("Empty title", "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "MetaData \"X-Image-Title\" is empty"})
		return
	}

	rawImage, err := io.ReadAll(c.Request.Body)
	if err != nil {
		ctrl.logger.Warn("Failed to read image bytes", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: invalid image in request body"})
		return
	}

	// rawImage, err := c.FormFile("image")
	// if err != nil {
	// 	ctrl.logger.Warn("Failed to read form file", logger.Err(err), "op", op)
	// 	ctrl.responce(c, http.StatusBadRequest, gin.H{"message": "Image not found in form-data"})
	// 	return
	// }

	// file, err := rawImage.Open()
	// if err != nil {
	// 	ctrl.logger.Warn("Failed to open file form multipart.FileHeader", logger.Err(err), "op", op)
	// 	ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Image is not valid"})
	// 	return
	// }
	// defer file.Close()

	// imageByte, err := getImageBuffer(file)
	// if err != nil {
	// 	ctrl.logger.Warn("Failed convert image to array of bytes", logger.Err(err), "op", op)
	// 	ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Image is not valid"})
	// 	return
	// }

	image := domain.Image{
		Title: title,
		Data:  rawImage,
	}

	if err := ctrl.service.Create(c.Request.Context(), &image); err != nil {
		if errors.Is(err, crud_errors.ErrImageCorruption) {
			ctrl.logger.Warn("Image data is corrupted", logger.Err(err), "op", op)
			ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid payload or image is corrapted"})
			return
		}

		ctrl.logger.Error("Failed to add image", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	resp := mapper.ImageDomainToImageResponse(image)
	ctrl.responce(c, http.StatusCreated, resp)
}

// GetAllImages godoc
//
//	@Summary		Get all images
//	@Description	The endpoint for retrieve all registered images in system
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	true	"limit get images"
//	@Param			offset	query		int	true	"offset get images"
//	@Success		200		{array}		dto.Image
//	@Failure		400		{object}	domain.Error
//	@Failure		404		{object}	domain.Error
//	@Failure		500		{object}	domain.Error
//	@Router			/api/v1/images [get]
func (ctrl *ImageController) GetAll(c *gin.Context) {
	op := "controllers.imageController.GetAll"
	limit, err := strconv.Atoi(c.DefaultQuery("limit", defaultLimit))
	if err != nil {
		ctrl.logger.Warn("Failed convert limit value", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: limit is not valid"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", defaultOffset))
	if err != nil {
		ctrl.logger.Warn("Failed convert offset value", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: offset is not valid"})
		return
	}

	images, err := ctrl.service.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		if errors.Is(err, crud_errors.ErrInvalidParam) {
			ctrl.logger.Warn("Invalid limit or offset parameter", "op", op)
			ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: limit cannot be less or equal 0, offset cannot be less than 0"})
			return
		}

		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("No content", logger.Err(err), "op", op)
			ctrl.responce(c, http.StatusNoContent, gin.H{"massage": "404: no data is contains"})
			return
		}

		ctrl.logger.Error("Failed to retrive images", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	output := make([]dto.ImageResponse, len(images))

	for i, image := range images {
		dto := mapper.ImageDomainToImageResponse(image)
		output[i] = dto
	}

	ctrl.logger.Debug("Retrieved all image's", "limit", limit, "offset", offset, "op", op)
	ctrl.responce(c, http.StatusOK, output)
}

// GetImageByID godoc
//
//	@Summary		Get all images
//	@Description	The endpoint for retrieve registered image in system by id
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uuid.UUID	true	"Image ID"
//	@Success		200	{object}	dto.Image
//	@Failure		400	{object}	domain.Error
//	@Failure		404	{object}	domain.Error
//	@Failure		500	{object}	domain.Error
//	@Router			/api/v1/images/{id} [get]
func (ctrl *ImageController) GetById(c *gin.Context) {
	op := "controllers.imageController.GetById"
	rawId := c.Param("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Warn("The received identifier is invalid", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: id is not valid"})
		return
	}

	image, err := ctrl.service.GetById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("image not found", "op", op)
			ctrl.responce(c, http.StatusNotFound, gin.H{"massage": "404: image not found"})
			return
		}

		ctrl.logger.Error("Failed to get data from database", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	output := mapper.ImageDomainToImageResponse(*image)
	ctrl.logger.Debug("Image retrieved", "id", id, "op", op)
	ctrl.responce(c, http.StatusOK, output)
}

// UpdateImage godoc
//
//	@Summary		Update image
//	@Description	The endpoint for updating image data by ID to a new image given by the user
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			id		path	uuid.UUID		true	"Image ID"
//	@Param			image	body	domain.Image	true	"New image data"
//	@Success		200
//	@Failure		400	{object}	domain.Error
//	@Failure		404	{object}	domain.Error
//	@Failure		500	{object}	domain.Error
//	@Router			/api/v1/images/{id} [patch]
func (ctrl *ImageController) Update(c *gin.Context) {
	op := "controllers.imageController.Update"

	if c.ContentType() != contentTypeOctetStream {
		ctrl.logger.Warn("Invalid content-type", "got", c.ContentType(), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"message": "Expected multipart/form-data"})
		return
	}

	rawId := c.Param("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Warn("The received identifier is invalid", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: id is not valid"})
		return
	}

	title := c.Request.Header.Get(headerXImageTitle)
	rawImage, err := io.ReadAll(c.Request.Body)
	if err != nil {
		ctrl.logger.Warn("Failed to read image bytes", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: invalid image in request body"})
		return
	}

	if len(rawImage) == 0 {
		ctrl.logger.Warn("Image not recived", "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: empty image data"})
		return
	}

	// rawImage, err := c.FormFile("image")
	// if err != nil {
	// 	ctrl.logger.Warn("Failed to read form file", logger.Err(err), "op", op)
	// 	ctrl.responce(c, http.StatusBadRequest, gin.H{"message": "Image not found in form-data"})
	// 	return
	// }

	// file, err := rawImage.Open()
	// if err != nil {
	// 	ctrl.logger.Warn("Failed to open file form multipart.FileHeader", logger.Err(err), "op", op)
	// 	ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Image is not valid"})
	// 	return
	// }
	// defer file.Close()

	// imageByte, err := getImageBuffer(file)
	// if err != nil {
	// 	ctrl.logger.Warn("Failed convert image to array of bytes", logger.Err(err), "op", op)
	// 	ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Image is not valid"})
	// 	return
	// }

	image := domain.Image{
		Id:    id,
		Title: title,
		Data:  rawImage,
	}

	if err := ctrl.service.Update(c.Request.Context(), &image); err != nil {
		if errors.Is(err, crud_errors.ErrImageCorruption) {
			ctrl.logger.Warn("Invalid image data", logger.Err(err), "op", op)
			ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: invalid image or image is corrupted"})
			return
		}

		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("Image not found", logger.Err(err), "op", op)
			ctrl.responce(c, http.StatusNotFound, gin.H{"massage": "404: image not found for update"})
			return
		}

		ctrl.logger.Error("Failed to update image", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	ctrl.logger.Debug("Image is updated", "id", id, "op", op)
	c.Status(http.StatusOK)
}

// DeleteImage godoc
//
//	@Summary		Delete image
//	@Description	The endpoint for deleting image data by ID
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			id	path	uuid.UUID	true	"Image ID"
//	@Success		204
//	@Failure		400	{object}	domain.Error
//	@Failure		404	{object}	domain.Error
//	@Failure		500	{object}	domain.Error
//	@Router			/api/v1/images/{id} [delete]
func (ctrl *ImageController) Delete(c *gin.Context) {
	op := "controllers.imageController.Delete"
	rawId := c.Param("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Warn("The received identifier is invalid", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalud request payload: id is not valid"})
		return
	}

	if err := ctrl.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("No content for this id", "id", id, "op", op)
			c.Status(http.StatusNoContent)
			return
		}

		ctrl.logger.Error("Failed to delete image from database", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	ctrl.logger.Debug("Image deleted", "id", id, "op", op)
	c.Status(http.StatusNoContent)
}
