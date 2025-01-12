package controllers

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientsController struct {
	service services.ClientsService
}

func NewClientController(clientService services.ClientsService) *ClientsController {
	return &ClientsController{
		service: clientService,
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.service.AddClient(c.Request.Context(), &client); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Client added successfully"})
}

// Get /api/v1/clients
// Retrieve all clients
// 200:
// 400:
// 500:
func (ctrl *ClientsController) GetAllClients(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "10"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client, err := ctrl.service.GetAllClients(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

// Get /api/v1/clients/search
// Search client by name and surname
// 200
// 400
// 404
func (ctrl *ClientsController) GetClientById(c *gin.Context) {
	name := c.Param("client_name")
	surname := c.Param("client_surname")

	client, err := ctrl.service.GetClientById(c.Request.Context(), name, surname)
	if client == nil && err == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

// Patch /api/v1/clients/:id/address
// Change a address id parameter by a given new id parameter
// 200
// 400
// 500
func (ctrl *ClientsController) ChangeAddressIdParameter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedFields model.UpdateAddressID

	if err := c.ShouldBindJSON(&updatedFields); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.service.ChangeAddressParameter(c.Request.Context(), id, int(updatedFields.AddressId)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address ID"})
		return
	}

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.service.DeleteClientById(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Address deleted successfully"})
}
