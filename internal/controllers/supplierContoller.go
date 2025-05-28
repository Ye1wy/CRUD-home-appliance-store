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

type supplierService interface {
	Create(ctx context.Context, supplier *domain.Supplier) error
	GetAll(ctx context.Context, limit, offset int) ([]domain.Supplier, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Supplier, error)
	UpdateAddress(ctx context.Context, id uuid.UUID, address *domain.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SupplierController struct {
	*BaseController
	service supplierService
}

func NewSupplierContoller(service supplierService, logger *logger.Logger) *SupplierController {
	controller := NewBaseContorller(logger)
	logger.Debug("Supplier controller is created")
	return &SupplierController{
		BaseController: controller,
		service:        service,
	}
}

func (ctrl *SupplierController) Create(c *gin.Context) {
	op := "controllers.supplierController.Create"
	var input dto.Supplier

	if err := c.ShouldBind(&input); err != nil {
		ctrl.logger.Warn("Failed to bind JSON for Create", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, err)
		return
	}

	supplier := mapper.SupplierToDomain(input)

	err := ctrl.service.Create(c, &supplier)
	if err != nil {
		ctrl.logger.Error("Failed create supplier", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, err)
		return
	}

	ctrl.logger.Debug("New supplier created", "op", op)
	c.Status(http.StatusCreated)
}

func (ctrl *SupplierController) GetAll(c *gin.Context) {
	op := "controllers.supplierController.GetAll"
	offset, err := strconv.Atoi(c.DefaultQuery("offset", defaultOffset))
	if err != nil {
		ctrl.logger.Warn("Failed convert offset value", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", defaultLimit))
	if err != nil {
		ctrl.logger.Warn("Failed convert limit value", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	supplier, err := ctrl.service.GetAll(c, limit, offset)
	if err != nil {
		if errors.Is(err, crud_errors.ErrInvalidParam) {
			ctrl.logger.Warn("Invalid value limit or offset", logger.Err(err), "op", op)
			ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Warn("No supplier data", "op", op)
			ctrl.responce(c, http.StatusNotFound, gin.H{"massage": "No data is contains"})
			return
		}

		ctrl.logger.Error("Failed retrieved data", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"massage": "Server is busy"})
		return
	}

	output := make([]dto.Supplier, len(supplier), cap(supplier))

	for i, supplier := range supplier {
		output[i] = mapper.SupplierToDTO(supplier)
	}

	ctrl.logger.Debug("All data is retrieved", "op", op)
	ctrl.responce(c, http.StatusOK, output)
}

func (ctrl *SupplierController) GetById(c *gin.Context) {
	op := "controllers.SupplierController.GetById"
	rawId := c.Param("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Warn("Invalid id payload", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	supplier, err := ctrl.service.GetById(c, id)
	if err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Warn("Supplier not found", "op", op)
			ctrl.responce(c, http.StatusNotFound, gin.H{"massage": "supplier not found"})
			return
		}

		ctrl.logger.Error("Failed to get supplier with id", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"massage": "Server is busy"})
		return
	}

	output := mapper.SupplierToDTO(*supplier)
	ctrl.logger.Debug("Data retrieved", "op", op)
	ctrl.responce(c, http.StatusOK, output)
}

func (ctrl *SupplierController) UpdateAddress(c *gin.Context) {
	op := "controllers.supplierController.UpdateAddress"
	rawId := c.Param("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Warn("Invalid id payload", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var input dto.Address

	if err := c.ShouldBind(&input); err != nil {
		ctrl.logger.Error("Failed to bind JSON for ChangeAddressIdParameter", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Invalid request payload"})
		return
	}

	address := mapper.AddressToDomain(input)

	if err := ctrl.service.UpdateAddress(c, id, &address); err != nil {
		ctrl.logger.Error("Failed to update address ID", "SupplierId", id, logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Failed to update address ID"})
		return
	}

	ctrl.logger.Info("Address ID updated successfully", "op", op)
	c.Status(http.StatusOK)
}

func (ctrl *SupplierController) Delete(c *gin.Context) {
	op := "controllers.supplierController.Delete"
	rawId := c.Param("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Warn("Invalid id payload", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ctrl.service.Delete(c.Request.Context(), id); err != nil {
		ctrl.logger.Error("Failed delete supplier by id", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"massage": "server is busy"})
		return
	}

	ctrl.logger.Debug("Successfuly deleted", "Supplier id:", id, "op", op)
	c.Status(http.StatusNoContent)
}
