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

// Create supplier godoc
//
//	@Summary		Create supplier
//	@Description	Supplier created from JSON or XML, for create endpoint required: name, phone_number, country, city, street
//	@Tags			suppliers
//	@Accept			json/xml
//	@Produce		json/xml
//	@Param			name			path	string	true	"Supplier name"
//	@Param			phone_number	path	string	true	"Supplier phone number"
//	@Param			country			path	string	true	"Supplier location country"
//	@Param			city			path	string	true	"Supplier location city"
//	@Param			street			path	string	true	"Supplier location street"
//	@Success		201				{object}
//	@Failure		400				{object}	domain.Error
//	@Failure		404				{object}	domain.Error
//	@Failure		500				{object}	domain.Error
//	@Router			/api/v1/suppliers [post]
func (ctrl *SupplierController) Create(c *gin.Context) {
	op := "controllers.supplierController.Create"
	var input dto.Supplier

	if err := c.ShouldBind(&input); err != nil {
		ctrl.logger.Warn("Failed to bind JSON/XML for create", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: invalid data received"})
		return
	}

	ctrl.logger.Debug("check data", "data", input)

	supplier := mapper.SupplierToDomain(input)

	if err := ctrl.service.Create(c, &supplier); err != nil {
		ctrl.logger.Error("Failed create supplier", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	ctrl.logger.Debug("Supplier created", "op", op)
	c.Status(http.StatusCreated)
}

// Get all supplier godoc
//
//	@Summary		Get all supplier
//	@Description	That endpoint retrieve all registered supplier in system
//	@Tags			suppliers
//	@Accept			json/xml
//	@Produce		json/xml
//	@Success		200	{object}	[]dto.Supplier
//	@Failure		400	{object}	domain.Error
//	@Failure		404	{object}	domain.Error
//	@Failure		500	{object}	domain.Error
//	@Router			/api/v1/suppliers [get]
func (ctrl *SupplierController) GetAll(c *gin.Context) {
	op := "controllers.supplierController.GetAll"

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

	supplier, err := ctrl.service.GetAll(c, limit, offset)
	if err != nil {
		if errors.Is(err, crud_errors.ErrInvalidParam) {
			ctrl.logger.Warn("Invalid value limit or offset", logger.Err(err), "op", op)
			ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: limit cannot be less or equal 0, offset cannot be less than 0"})
			return
		}

		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("No supplier data", "op", op)
			ctrl.responce(c, http.StatusNotFound, gin.H{"massage": "404: no data is contains"})
			return
		}

		ctrl.logger.Error("Failed retrieved data", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	output := make([]dto.Supplier, len(supplier))

	for i, supplier := range supplier {
		output[i] = mapper.SupplierToDTO(supplier)
	}

	ctrl.logger.Debug("Retrieved all supplier's", "limit", limit, "offset", offset, "op", op)
	ctrl.responce(c, http.StatusOK, output)
}

// Get supplier godoc
//
//	@Summary		Get supplier by ID
//	@Description	That endpoint retrieve registered supplier in system by ID
//	@Tags			suppliers
//	@Accept			json/xml
//	@Produce		json/xml
//	@Success		200	{object}	dto.Supplier
//	@Failure		400	{object}	domain.Error
//	@Failure		404	{object}	domain.Error
//	@Failure		500	{object}	domain.Error
//	@Router			/api/v1/suppliers/:id [get]
func (ctrl *SupplierController) GetById(c *gin.Context) {
	op := "controllers.SupplierController.GetById"
	rawId := c.Param("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Warn("The received identifier is invalid", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalud request payload: id is not valid"})
		return
	}

	supplier, err := ctrl.service.GetById(c, id)
	if err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("Supplier not found", "op", op)
			ctrl.responce(c, http.StatusNotFound, gin.H{"massage": "404: supplier not found"})
			return
		}

		ctrl.logger.Error("Failed to get supplier with id", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	output := mapper.SupplierToDTO(*supplier)
	ctrl.logger.Debug("Supplier retrieved", "id", id, "op", op)
	ctrl.responce(c, http.StatusOK, output)
}

// Update supplier godoc
//
//	@Summary		Update supplier by ID
//	@Description	That endpoint update supplier data (change address on supplier)
//	@Tags			suppliers
//	@Accept			json/xml
//	@Produce		json/xml
//	@Success		200	{object}
//	@Failure		400	{object}	domain.Error
//	@Failure		404	{object}	domain.Error
//	@Failure		500	{object}	domain.Error
//	@Router			/api/v1/suppliers/:id [patch]
func (ctrl *SupplierController) UpdateAddress(c *gin.Context) {
	op := "controllers.supplierController.UpdateAddress"
	rawId := c.Param("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Warn("The received identifier is invalid", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalud request payload: id is not valid"})
		return
	}

	var input dto.Address

	if err := c.ShouldBind(&input); err != nil {
		ctrl.logger.Warn("Failed to bind JSON/XML for update address", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: invalid data received"})
		return
	}

	address := mapper.AddressToDomain(input)

	if err := ctrl.service.UpdateAddress(c, id, &address); err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("Supplier not found", logger.Err(err), "op", op)
			ctrl.responce(c, http.StatusNotFound, gin.H{"massage": "404: supplier not found for update"})
			return
		}

		ctrl.logger.Error("Failed to update address ID", "id", id, logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	ctrl.logger.Debug("Supplier updated", "id", id, "op", op)
	c.Status(http.StatusOK)
}

// Delete supplier godoc
//
//	@Summary		Delete supplier by ID
//	@Description	That endpoint delete supplier data by id
//	@Tags			suppliers
//	@Accept			json/xml
//	@Produce		json/xml
//	@Success		204	{object}
//	@Failure		400	{object}	domain.Error
//	@Failure		404	{object}	domain.Error
//	@Failure		500	{object}	domain.Error
//	@Router			/api/v1/suppliers/:id [delete]
func (ctrl *SupplierController) Delete(c *gin.Context) {
	op := "controllers.supplierController.Delete"
	rawId := c.Param("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Warn("The received identifier is invalid", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalud request payload: id is not valid"})
		return
	}

	if err := ctrl.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("Supplier not found", "op", op)
			c.Status(http.StatusNoContent)
			return
		}

		ctrl.logger.Error("Failed delete supplier by id", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	ctrl.logger.Debug("Supplier deleted", "id", id, "op", op)
	c.Status(http.StatusNoContent)
}
