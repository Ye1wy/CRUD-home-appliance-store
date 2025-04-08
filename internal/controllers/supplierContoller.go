package controllers

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SupplierController struct {
	*BaseController
	service services.SupplierServiceInterface
}

func NewSupplierContoller(service services.SupplierServiceInterface, logger *logger.Logger) *SupplierController {
	controller := NewBaseContorller(logger)
	logger.Debug("Supplier controller is created")
	return &SupplierController{
		BaseController: controller,
		service:        service,
	}
}

func (ctrl *SupplierController) Create(c *gin.Context) {
	op := "controllers.supplierController.Create"
	var dto dto.SupplierDTO

	if err := ctrl.mapping(c, dto); err != nil {
		ctrl.logger.Error("Failed to bind JSON for Create", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, err)
		return
	}

	supplier, err := ctrl.service.Create(c, &dto)
	if err != nil {
		ctrl.logger.Error("Failed create supplier", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, err)
		return
	}

	ctrl.logger.Debug("New supplier created", "op", op)
	ctrl.responce(c, http.StatusCreated, supplier)
}

func (ctrl *SupplierController) GetAll(c *gin.Context) {
	op := "controllers.supplierController.GetAll"
	offset, err := strconv.Atoi(c.DefaultQuery("offset", defaultOffset))
	if err != nil {
		ctrl.logger.Error("Failed convert offset value", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, err)
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", defaultLimit))
	if err != nil {
		ctrl.logger.Error("Failed convert limit value", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, err)
		return
	}

	data, err := ctrl.service.GetAll(c, limit, offset)
	if err != nil {
		ctrl.logger.Error("Failed retrieved data", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, err)
		return
	}

	ctrl.logger.Debug("All data is retrieved", "op", op)
	ctrl.responce(c, http.StatusOK, data)
}

func (ctrl *SupplierController) GetById(c *gin.Context) {
	op := "controllers.SupplierController.GetById"
	id := c.Param("id")

	supplierDTO, err := ctrl.service.GetById(c, id)
	if supplierDTO == nil && err == nil {
		ctrl.logger.Warn("Supplier not found", "op", op)
		ctrl.responce(c, http.StatusNotFound, gin.H{"msg": "supplier not found"})
		return
	}

	if err != nil {
		ctrl.logger.Error("Failed to get supplier with id", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, err)
		return
	}

	ctrl.logger.Debug("Data retrieved", "op", op)
	ctrl.responce(c, http.StatusOK, supplierDTO)
}

func (ctrl *SupplierController) UpdateAddress(c *gin.Context) {
	op := "controllers.supplierController.UpdateAddress"
	id := c.Param("id")
	var address dto.UpdateAddressDTO

	if err := ctrl.mapping(c, address); err != nil {
		ctrl.logger.Error("Failed to bind JSON for ChangeAddressIdParameter", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ctrl.service.UpdateAddress(c, id, address.AddressID); err != nil {
		ctrl.logger.Error("Failed to update address ID", "SupplierId", id, logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Failed to update address ID"})
		return
	}

	ctrl.logger.Info("Address ID updated successfully", "op", op)
	ctrl.responce(c, http.StatusOK, gin.H{"Msg": "Object updated"})
}

func (ctrl *SupplierController) Delete(c *gin.Context) {
	op := "controllers.supplierController.Delete"
	id := c.Param("id")

	if err := ctrl.service.Delete(c.Request.Context(), id); err != nil {
		ctrl.logger.Error("Failed delete supplier by id", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, err)
		return
	}

	ctrl.logger.Debug("Successfuly deleted", "SupplierID", id, "op", op)
	c.Status(http.StatusNoContent)
}
