package controllers

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductsController struct {
	service services.ProductService
	logger  *logger.Logger
}

func NewProductController(service services.ProductService, log *logger.Logger) *ProductsController {
	return &ProductsController{
		service: service,
		logger:  log,
	}
}

// Post /api/v1/products
// Add product
// 201
// 400
// 500
func (ctrl *ProductsController) AddProduct(c *gin.Context) {
	op := "controllers.product.AddProduct"

	var productDTO dto.ProductDTO
	if err := c.ShouldBindJSON(&productDTO); err != nil {
		ctrl.logger.Error("Failed to bind JSON for AddProduct", logger.Err(err), "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	product, err := ctrl.service.AddProduct(c.Request.Context(), &productDTO)
	if err != nil {
		ctrl.logger.Error("Failed to add product: ", logger.Err(err), "op", op)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product"})
		return
	}

	ctrl.logger.Info("Product added successfully", "productID", product.Id, "op", op)
	c.JSON(http.StatusCreated, product)
}

// Get /api/v1/products?limit=&offset=
// Retrieve all product
func (ctrl *ProductsController) GetAllProduct(c *gin.Context) {
	op := "controllers.product.GetAllProduct"

	limit, err := strconv.Atoi(c.DefaultQuery("limit", defaultLimit))
	if err != nil {
		ctrl.logger.Error("Invalid limit parameter", logger.Err(err), "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", defaultOffset))
	if err != nil {
		ctrl.logger.Error("Invalid offset parameter", logger.Err(err), "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	client, err := ctrl.service.GetAllProducts(c.Request.Context(), limit, offset)
	if err != nil {
		ctrl.logger.Error("Failed to retrieve product", logger.Err(err), "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve product"})
	}

	ctrl.logger.Info("retrieved all products", "limit", limit, "offset", offset)
	c.JSON(http.StatusOK, client)
}

// Get /api/v1/products/:id
// Get product by id
func (ctrl *ProductsController) GetProductById(c *gin.Context) {
	op := "controllers.product.GetProductById"

	id := c.Param("id")

	productDTO, err := ctrl.service.GetProductById(c.Request.Context(), id)

	if productDTO == nil && err == nil {
		ctrl.logger.Warn("Product not found", "op", op)
		c.JSON(http.StatusNotFound, gin.H{"warning": "Product not found"})
		return
	}

	if err != nil {
		ctrl.logger.Error("Error while searching product", "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id cannot be empty or id is invalid"})
		return
	}

	ctrl.logger.Info("Product retrieved successfully", "id", id)
	c.JSON(http.StatusOK, productDTO)
}

// Patch /api/v1/products/:id/decrease=?
// Decrease a parameter by a given value
func (ctrl *ProductsController) DecreaseParameter(c *gin.Context) {
	op := "controllers.product.DecreaseParameter"

	id := c.Param("id")
	decreaseValue, err := strconv.Atoi(c.Query("decrease"))
	if err != nil {
		ctrl.logger.Error("Error in convert string from query to int", "op", op)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if err := ctrl.service.DecreaseStock(c.Request.Context(), id, decreaseValue); err != nil {
		ctrl.logger.Error("Error in decrease stock", "op", op)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctrl.logger.Info("Stock in product is decreased successfully", "op", op)
	c.Status(http.StatusOK)
}

// Delete /api/v1/products/:id
// Delete product by identificator
func (ctrl *ProductsController) DeleteProductById(c *gin.Context) {
	op := "controllers.product.DeleteProductById"

	id := c.Param("id")

	if err := ctrl.service.DeleteProductById(c.Request.Context(), id); err != nil {
		ctrl.logger.Error("Failed to delete product", logger.Err(err), "op", op)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	ctrl.logger.Info("Product deleted successfully", "productID", id, "op", op)
	c.Status(http.StatusNoContent)
}
