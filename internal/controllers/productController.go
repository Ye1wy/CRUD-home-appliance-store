package controllers

import (
	crud_errors "CRUD-HOME-APPLIANCE-STORE/internal/errors"
	"CRUD-HOME-APPLIANCE-STORE/internal/mapper"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type productService interface {
	Create(ctx context.Context, product *domain.Product) error
	GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	Update(ctx context.Context, id uuid.UUID, decrease int) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProductController struct {
	*BaseController
	service productService
}

func NewProductController(service productService, logger *logger.Logger) *ProductController {
	controller := NewBaseContorller(logger)
	logger.Debug("Product controller is created")
	return &ProductController{
		BaseController: controller,
		service:        service,
	}
}

// Create Product godoc
//
//		@Summary		Create product
//		@Description	Product created from JSON or XML, for create endpoint required: name, category, price, available_stock, supplier_name, supplier_phone_number, image
//		@Tags			products
//		@Accept			json/xml
//		@Produce		json/xml
//		@Param			name	path	string true "Product name"
//	 	@Param			category path	string true "Product categoty"
//	 	@Param			price path 	float64 true "Product price"
//	 	@Param			available_stock path		int64 true "Prodcut available stock"
//	 	@Param			supplier_name path string true "Product supplier name"
//	 	@Param			supplier_phone_number path string true "Product supplier phone number"
//	 	@Param			image path []byte true "Product image"
//		@Success		201	{object}
//		@Failure		400	{object}	domain.Error
//		@Failure		404	{object}	domain.Error
//		@Failure		500	{object}	domain.Error
//		@Router			/api/v1/products [post]
func (ctrl *ProductController) Create(c *gin.Context) {
	op := "controllers.productController.Create"
	var input dto.Product

	if err := c.ShouldBind(&input); err != nil {
		ctrl.logger.Warn("Failed to bind JSON/XML for AddProduct", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: invalid data received"})
		return
	}

	product := mapper.ProductToDomain(input)

	if err := ctrl.service.Create(c.Request.Context(), &product); err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("Invalid supplier data with create product", "op", op)
			ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid supplier data"})
			return
		}

		ctrl.logger.Error("Failed to add client", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	ctrl.logger.Debug("Product created", "id", product.Id, "op", op)
	c.Status(http.StatusCreated)
}

// Get All Product godoc
//
//		@Summary		Get all product
//	@Description	The endpoint for retrieve all registered product in system
//		@Tags			products
//		@Accept			json/xml
//		@Produce		json/xml
//		@Success		200	{object}	[]dto.Product
//		@Failure		400	{object}	domain.Error
//		@Failure		404	{object}	domain.Error
//		@Failure		500	{object}	domain.Error
//		@Router			/api/v1/products [get]
func (ctrl *ProductController) GetAll(c *gin.Context) {
	op := "controllers.productController.GetAll"
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

	product, err := ctrl.service.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		if errors.Is(err, crud_errors.ErrInvalidParam) {
			ctrl.logger.Warn("Invalid limit or offset parameter", "op", op)
			ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: limit cannot be less or equal 0, offset cannot be less than 0"})
			return
		}

		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("Product's not found: no content", "op", op)
			ctrl.responce(c, http.StatusNotFound, gin.H{"massage": "404: no data is contains"})
			return
		}

		ctrl.logger.Error("Failed to retrieve product", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	output := make([]dto.Product, len(product), cap(product))

	for i, item := range product {
		dto := mapper.ProductToDTO(item)
		output[i] = dto
	}

	ctrl.logger.Debug("Retrieved all products", "limit", limit, "offset", offset)
	ctrl.responce(c, http.StatusOK, output)
}

// Get Product godoc
//
//		@Summary		Get product by id
//	@Description	The endpoint for retrieve registered product in system by id
//		@Tags			products
//		@Accept			json/xml
//		@Produce		json/xml
//		@Success		200	{object}	dto.Product
//		@Failure		400	{object}	domain.Error
//		@Failure		404	{object}	domain.Error
//		@Failure		500	{object}	domain.Error
//		@Router			/api/v1/products/:id [get]
func (ctrl *ProductController) GetById(c *gin.Context) {
	op := "controllers.productController.GetById"
	rawId := c.Param("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Warn("The received identifier is invalid", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalud request payload: id is not valid"})
		return
	}

	product, err := ctrl.service.GetById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("Product not found", "op", op)
			ctrl.responce(c, http.StatusNotFound, gin.H{"massage": "404: product not found"})
			return
		}

		ctrl.logger.Error("Failed to get prodcut with id", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	output := mapper.ProductToDTO(*product)
	ctrl.logger.Debug("Product retrieved", "id", id, "op", op)
	ctrl.responce(c, http.StatusOK, output)
}

// Update Product godoc
//
//		@Summary		Update product by ID
//	@Description	The endpoint for updating product data (avaliable stock) by ID to a decrease avalible stock
//		@Tags			products
//		@Accept			json/xml
//		@Produce		json/xml
//		@Success		201	{object}
//		@Failure		400	{object}	domain.Error
//		@Failure		404	{object}	domain.Error
//		@Failure		500	{object}	domain.Error
//		@Router			/api/v1/products/:id?decrease= [patch]
func (ctrl *ProductController) Update(c *gin.Context) {
	op := "controllers.productController.DecreaseStock"
	rawId := c.Param("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Warn("The received identifier is invalid", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalud request payload: id is not valid"})
		return
	}

	value, err := strconv.Atoi(c.Query("decrease"))
	if err != nil {
		ctrl.logger.Warn("Failed convert decrease value", "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: decrease value is invalid"})
		return
	}

	if err := ctrl.service.Update(c.Request.Context(), id, value); err != nil {
		if errors.Is(err, crud_errors.ErrInvalidParam) {
			ctrl.logger.Warn("Invalid deacrease value", "value", value, "op", op)
			ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid payload: decrease value cannot be less than 0"})
			return
		}

		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("Product not found", logger.Err(err), "op", op)
			ctrl.responce(c, http.StatusNotFound, gin.H{"massage": "404: product not found for update"})
			return
		}

		ctrl.logger.Error("Failed to update avalilable stock", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	ctrl.logger.Debug("Product updated", "op", op)
	c.Status(http.StatusOK)
}

// Delete Product godoc
//
//		@Summary		Delete product by ID
//	@Description	The endpoint for deleting product data by ID
//		@Tags			products
//		@Accept			json/xml
//		@Produce		json/xml
//		@Success		204 {object}
//		@Failure		400	{object}	domain.Error
//		@Failure		404	{object}	domain.Error
//		@Failure		500	{object}	domain.Error
//		@Router			/api/v1/products/:id[delete]
func (ctrl *ProductController) Delete(c *gin.Context) {
	op := "controllers.productController.Delete"
	rawId := c.Param("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Warn("The received identifier is invalid", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalud request payload: id is not valid"})
		return
	}

	if err := ctrl.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("Product not found", "op", op)
			c.Status(http.StatusNoContent)
			return
		}

		ctrl.logger.Error("Failed to delete product by id", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	ctrl.logger.Debug("Product deleted", "id", id, "op", op)
	c.Status(http.StatusNoContent)
}
