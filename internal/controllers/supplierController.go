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

type SuppliersController struct {
	service services.SuppliersService
	logger  *slog.Logger
}

func NewSuppliersController(supplierService services.SuppliersService, log *slog.Logger) *SuppliersController {
	return &SuppliersController{
		service: supplierService,
		logger:  log,
	}
}

// Post /api/v1/suppliers
// Add supplier
// 201:
// 400
// 500
func (ctrl *SuppliersController) AddSupplier(c *gin.Context) {
	op := "controllers.supplier.addSupplier"

	var supplier model.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		ctrl.logger.Error("Failed to bind JSON for AddSupplier", logger.Err(err), "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ctrl.service.AddSupplier(c.Request.Context(), &supplier); err != nil {
		ctrl.logger.Error("Failed to add supplier: ", logger.Err(err), "op", op)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add supplier"})
		return
	}

	ctrl.logger.Info("Supplier added successfully", "supplierID", supplier.Id, "op", op)
	c.JSON(http.StatusOK, gin.H{"message": "Supplier added successfully"})
}

// Get /api/v1/suppliers
// Retrieve all suppliers
// 200:
// 400:
// 500:
func (ctrl *SuppliersController) GetAllSuppliers(c *gin.Context) {
	op := "controllers.supplier.getAllSuppliers"

	limit, err := strconv.Atoi(c.DefaultQuery("limit", defaultLimit))
	if err != nil {
		ctrl.logger.Error("Invalid limit parameter", logger.Err(err), "op", op)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", defaultOffset))
	if err != nil {
		ctrl.logger.Error("Invalid offset parameter", logger.Err(err), "op", op)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid offset parameter"})
		return
	}

	supplier, err := ctrl.service.GetAllSuppliers(c.Request.Context(), limit, offset)
	if err != nil {
		ctrl.logger.Error("Failed to retrieve supplier", logger.Err(err), "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve supplier"})
		return
	}

	ctrl.logger.Info("Retrieved all suppliers", "limit", limit, "offset", offset, "op", op)
	c.JSON(http.StatusOK, supplier)
}

// Get /api/v1/suppliers/:id
// Get supplier by id
// 200
// 400
// 500
func (ctrl *SuppliersController) GetSupplierById(c *gin.Context) {
	op := "controllers.supplier.getSupplierById"

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		ctrl.logger.Error("Failed to take parameter from query", logger.Err(err), "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to take parameter"})
		return
	}

	supplier, err := ctrl.service.GetSupplierById(c.Request.Context(), id)
	if supplier == nil && err == nil {
		ctrl.logger.Warn("Supplier not found", logger.Err(err), "op", op)
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}

	if err != nil {
		ctrl.logger.Error("Error while searching supplier ", logger.Err(err), "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to found supplier"})
		return
	}

	ctrl.logger.Info("Supplier retrieved successefully", "id", id, "op", op)
	c.JSON(http.StatusOK, supplier)
}

// Patch /api/v1/suppliers/:id/changeAddress
// Change a address id parameter by a given value
// 200
// 400
// 500
func (ctrl *SuppliersController) ChangeAddressParameter(c *gin.Context) {
	op := "controllers.supplier.changeAddressParameter"

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ctrl.logger.Error("Invalid supplier id parameter", logger.Err(err), "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
		return
	}

	var updatedFields model.UpdateAddressID

	if err := c.ShouldBindJSON(&updatedFields); err != nil {
		ctrl.logger.Error("Failed to bind JSON for ChangeAddressIdParameter", logger.Err(err), "op", op)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ctrl.service.ChangeAddressParameter(c.Request.Context(), id, int(updatedFields.AddressId)); err != nil {
		ctrl.logger.Error("Failed to update address ID", "supplierID", id, logger.Err(err), "op", op)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address ID"})
		return
	}

	ctrl.logger.Info("Address ID updated successfully", "supplierID", id, "newAddressID", updatedFields.AddressId, "op", op)
	c.JSON(http.StatusOK, gin.H{"status": "Address ID updated"})
}

// Delete /api/v1/suppliers/:id
// Delete supplier by identificator
func (ctrl *SuppliersController) DeleteSupplierById(c *gin.Context) {
	op := "controllers.supplier.deleteSupplierById"

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ctrl.logger.Error("Invalid supplier ID parameter for DeleteSupplierById", logger.Err(err), "op", op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
		return
	}

	if err := ctrl.service.DeleteSupplierById(c.Request.Context(), id); err != nil {
		ctrl.logger.Error("Failed to delete supplier", "supplierID", id, logger.Err(err), "op", op)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete supplier"})
		return
	}

	ctrl.logger.Info("Supplier deleted successfully", "supplierID", id, "op", op)
	c.Status(http.StatusNoContent)
}
