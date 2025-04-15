package controllers

// import (
// 	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
// 	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
// 	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
// 	"context"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// )

// type ProductServiceInterface interface {
// 	Create(ctx context.Context, product *domain.Product) error
// 	GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error)
// 	GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
// 	Update(ctx context.Context, id uuid.UUID, decrease int) error
// 	Delete(ctx context.Context, id uuid.UUID) error
// }

// type ProductController struct {
// 	*BaseController
// 	service ProductServiceInterface
// }

// func NewProductController(service ProductServiceInterface, logger *logger.Logger) *ProductController {
// 	controller := NewBaseContorller(logger)
// 	logger.Debug("Product controller is created")
// 	return &ProductController{
// 		BaseController: controller,
// 		service:        service,
// 	}
// }

// // Post /api/v1/products
// // Add product
// // 201
// // 400
// // 500
// func (ctrl *ProductController) Create(c *gin.Context) {
// 	op := "controllers.productController.Create"
// 	var productDTO dto.ProductDTO

// 	if err := ctrl.mapping(c, &productDTO); err != nil {
// 		ctrl.logger.Error("Failed to bind JSON/XML for AddProduct", logger.Err(err), "op", op)
// 		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 		return
// 	}

// 	product, err := ctrl.service.Create(c.Request.Context(), &productDTO)
// 	if err != nil {
// 		ctrl.logger.Error("Failed to add client: ", logger.Err(err), "op", op)
// 		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Failed to add product"})
// 		return
// 	}

// 	ctrl.logger.Debug("Product added successfully", "productID", product.Id, "op", op)
// 	ctrl.responce(c, http.StatusCreated, product)
// }

// // Get /api/v1/products?limit=&offset=
// // Retrieve all product
// func (ctrl *ProductController) GetAll(c *gin.Context) {
// 	op := "controllers.productController.GetAll"

// 	limit, err := strconv.Atoi(c.DefaultQuery("limit", defaultLimit))
// 	if err != nil {
// 		ctrl.logger.Error("Invalid limit parameter", logger.Err(err), "op", op)
// 		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Invalid payload"})
// 		return
// 	}

// 	offset, err := strconv.Atoi(c.DefaultQuery("offset", defaultOffset))
// 	if err != nil {
// 		ctrl.logger.Error("Invalid offset parameter", logger.Err(err), "op", op)
// 		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Invalid payload"})
// 		return
// 	}

// 	productDTOs, err := ctrl.service.GetAll(c.Request.Context(), limit, offset)
// 	if err != nil {
// 		ctrl.logger.Error("Failed to retrieve product", logger.Err(err), "op", op)
// 		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Failed to retrieve product"})
// 		return
// 	}

// 	ctrl.logger.Debug("retrieved all products", "limit", limit, "offset", offset)
// 	ctrl.responce(c, http.StatusOK, productDTOs)
// }

// // Get /api/v1/products/:id
// // Get product by id
// func (ctrl *ProductController) GetById(c *gin.Context) {
// 	op := "controllers.productController.GetById"
// 	id := c.Param("id")

// 	productDTO, err := ctrl.service.GetById(c.Request.Context(), id)
// 	if productDTO == nil && err == nil {
// 		ctrl.logger.Warn("Product not found", "op", op)
// 		ctrl.responce(c, http.StatusNotFound, gin.H{"warning": "product not found"})
// 		return
// 	}

// 	if err != nil {
// 		ctrl.logger.Error("Error while searching product", "op", op)
// 		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "id cannot be empty or id is invalid"})
// 		return
// 	}

// 	ctrl.logger.Debug("Product retrieved successfully", "id", id)
// 	ctrl.responce(c, http.StatusOK, productDTO)
// }

// // Patch /api/v1/products/:id/decrease=?
// // Decrease a parameter by a given value
// func (ctrl *ProductController) DecreaseStock(c *gin.Context) {
// 	op := "controllers.productController.DecreaseStock"
// 	id := c.Param("id")
// 	decreaseValue, err := strconv.Atoi(c.Query("decrease"))
// 	if err != nil {
// 		ctrl.logger.Error("Error in convert string from query to int", "op", op)
// 		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": err})
// 		return
// 	}

// 	if err := ctrl.service.Update(c.Request.Context(), id, decreaseValue); err != nil {
// 		ctrl.logger.Error("Error in decrease stock", "op", op)
// 		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": err})
// 		return
// 	}

// 	ctrl.logger.Debug("Stock in product is decreased successfully", "op", op)
// 	c.Status(http.StatusOK)
// }

// // Delete /api/v1/products/:id
// // Delete product by identificator
// func (ctrl *ProductController) Delete(c *gin.Context) {
// 	op := "controllers.productController.Delete"
// 	id := c.Param("id")

// 	if err := ctrl.service.Delete(c.Request.Context(), id); err != nil {
// 		ctrl.logger.Error("Failed to delete product", logger.Err(err), "op", op)
// 		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
// 		return
// 	}
// 	ctrl.logger.Debug("Product deleted successfully", "productID", id, "op", op)
// 	c.Status(http.StatusNoContent)
// }
