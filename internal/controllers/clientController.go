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

const (
	DefaultLimit  = "10"
	DefaultOffset = "0"
)

type ClientsController struct {
	service services.ClientsService
	logger  *slog.Logger
}

func NewClientController(clientService services.ClientsService, logger *slog.Logger) *ClientsController {
	return &ClientsController{
		service: clientService,
		logger:  logger,
	}
}

// Postapi/v1/clients
// Add client
// 201:
// 400
// 500
func (ctrl *ClientsController) AddClient(c *gin.Context) {
	var client model.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		ctrl.logger.Error("[ERROR] Failed to bind JSON for AddClient", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ctrl.service.AddClient(c.Request.Context(), &client); err != nil {
		ctrl.logger.Error("[ERROR] Failed to add client: ", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add client"})
		return
	}

	ctrl.logger.Info("[INFO] Client added successfully", "clientID", client.Id)
	c.JSON(http.StatusCreated, gin.H{"message": "Client added successfully"})
}

// Get /api/v1/clients
// Retrieve all clients
// 200:
// 400:
// 500:
func (ctrl *ClientsController) GetAllClients(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", DefaultLimit))
	if err != nil {
		ctrl.logger.Error("[ERROR] Invalid limit parameter", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", DefaultOffset))
	if err != nil {
		ctrl.logger.Error("[ERROR] Invalid offset parameter", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid offset parameter"})
		return
	}

	client, err := ctrl.service.GetAllClients(c.Request.Context(), limit, offset)
	if err != nil {
		ctrl.logger.Error("[ERROR] Failed to retrieve clients", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.logger.Info("[INFO] Retrieved all clients", "limit", limit, "offset", offset)
	c.JSON(http.StatusOK, client)
}

// Get /api/v1/clients/search
// Search client by name and surname
// 200
// 400
// 404
func (ctrl *ClientsController) GetClientById(c *gin.Context) {
	name := c.Query("client_name")
	surname := c.Query("client_surname")

	client, err := ctrl.service.GetClientById(c.Request.Context(), name, surname)
	if client == nil && err == nil {
		ctrl.logger.Warn("[WARN] Client not found", logger.Err(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	if err != nil {
		ctrl.logger.Error("[ERROR] Error while searching client ", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to found client"})
		return
	}

	ctrl.logger.Info("[INFO] Client retrieved successfully", "name", name, "surname", surname)
	c.JSON(http.StatusOK, client)
}

// Patch /api/v1/clients/:id/address
// Change a address id parameter by a given new id parameter
// 200
// 400
// 500
func (ctrl *ClientsController) ChangeAddressParameter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ctrl.logger.Error("[ERROR] Invalid client ID parameter", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	var updatedFields model.UpdateAddressID

	if err := c.ShouldBindJSON(&updatedFields); err != nil {
		ctrl.logger.Error("[ERROR] Failed to bind JSON for ChangeAddressIdParameter", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ctrl.service.ChangeAddressParameter(c.Request.Context(), id, int(updatedFields.AddressId)); err != nil {
		ctrl.logger.Error("[ERROR] Failed to update address ID", "clientID", id, logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address ID"})
		return
	}

	ctrl.logger.Info("[INFO] Address ID updated successfully", "clientID", id, "newAddressID", updatedFields.AddressId)
	c.JSON(http.StatusOK, gin.H{"status": "Address ID updated"})
}

// Delete /api/v1/clients/:id
// Delete client by identificator
// 204
// 400
// 500
func (ctrl *ClientsController) DeleteClientById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ctrl.logger.Error("[ERROR] Invalid client ID parameter for DeleteClientById", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	if err := ctrl.service.DeleteClientById(c.Request.Context(), id); err != nil {
		ctrl.logger.Error("Failed to delete client", "clientID", id, logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete client"})
		return
	}

	ctrl.logger.Info("Client deleted successfully", "clientID", id)
	c.Status(http.StatusNoContent)
}
