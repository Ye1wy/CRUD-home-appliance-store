package controllers

import (
	"github.com/gin-gonic/gin"
)

type SuppliersAPI struct {
}

// Post /api/v1/api/v1/suppliers
// Add supplier
func (api *SuppliersAPI) AddSupplier(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Patch /api/v1/api/v1/suppliers/:id/changeAddress
// Change a address id parameter by a given value
func (api *SuppliersAPI) ChangeAddressParametr(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Delete /api/v1/api/v1/suppliers/:id
// Delete supplier by identificator
func (api *SuppliersAPI) DeleteSupplierById(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Get /api/v1/api/v1/suppliers
// Retrieve all suppliers
func (api *SuppliersAPI) GetAllSuppliers(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}

// Get /api/v1/api/v1/suppliers/:id
// Get supplier by id
func (api *SuppliersAPI) SearchSupplierById(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}
