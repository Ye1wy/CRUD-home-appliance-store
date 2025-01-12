package controllers

import (
	"github.com/gin-gonic/gin"
)

type ClientsAPI struct {
}

// Post /api/v1/api/v1/clients
// Add client
func (api *ClientsAPI) AddClient(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Patch /api/v1/api/v1/clients/:id/address
// Change a address id parameter by a given new id parameter
func (api *ClientsAPI) ChangeAddressIdParameter(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Delete /api/v1/api/v1/clients/:id
// Delete client by identificator
func (api *ClientsAPI) DeleteClientById(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Get /api/v1/api/v1/clients
// Retrieve all clients
func (api *ClientsAPI) GetAllClients(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Get /api/v1/api/v1/clients/search
// Search client by name and surname
func (api *ClientsAPI) SearchClientByNameAndSurname(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}
