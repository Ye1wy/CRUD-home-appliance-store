package controllers

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/mapper"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
	psgrep "CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClientsServiceInterface interface {
	Create(ctx context.Context, client *domain.Client) error
	GetAll(ctx context.Context, limit, offset int) ([]domain.Client, error)
	GetByNameAndSurname(ctx context.Context, name, surname string) ([]domain.Client, error)
	UpdateAddress(ctx context.Context, object *domain.Client) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ClientController struct {
	*BaseController
	service ClientsServiceInterface
}

func NewClientsController(service ClientsServiceInterface, logger *logger.Logger) *ClientController {
	controller := NewBaseContorller(logger)
	return &ClientController{
		BaseController: controller,
		service:        service,
	}
}

// Create Client godoc
//
//		@Summary		Create client
//		@Description	Client created from JSON or XML, for create endpoint required: name, surname, birthday, gender, address_id
//		@Tags			clients
//		@Accept			json/xml
//		@Produce		json/xml
//		@Param			name	path	string true "Client name"
//	 @Param			surname path	string true "Client surname"
//	 @Param			birthday path 	string true "Client birthday"
//	 @Param			gender path		string true "Client gender"
//	 @Param			address_id path string true uuid.UUID "Client living address"
//		@Success		200	{object}	dto.ClientDTO
//		@Failure		400	{object}	dto.ErrorDTO
//		@Failure		404	{object}	dto.ErrorDTO
//		@Failure		500	{object}	dto.ErrorDTO
//		@Router			/api/v1/clients [post]
func (ctrl *ClientController) Create(c *gin.Context) {
	op := "controllers.clientController.Create"
	var clientDTO dto.ClientDTO

	if err := ctrl.mapping(c, &clientDTO); err != nil {
		ctrl.logger.Error("Failed to bind JSON for Create", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	client, err := mapper.ClientToDomain(&clientDTO)
	if err != nil {
		ctrl.logger.Error("Failed mapping dto to domain", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Error from server"})
		return
	}

	if err := ctrl.service.Create(c.Request.Context(), &client); err != nil {
		ctrl.logger.Error("Failed to add client: ", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Failed to add client"})
		return
	}

	ctrl.logger.Info("Client added successfully", "op", op)
	ctrl.responce(c, http.StatusCreated, client)
}

// Get /api/v1/clients
// Retrieve all clients
// 200:
// 400:
// 500:
func (ctrl *ClientController) GetAll(c *gin.Context) {
	op := "controllers.clientController.getAll"
	limit, err := strconv.Atoi(c.DefaultQuery("limit", defaultLimit))
	if err != nil {
		ctrl.logger.Error("Invalid limit parameter", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Invalid payload"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", defaultOffset))
	if err != nil {
		ctrl.logger.Error("Invalid offset parameter", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Invalid payload"})
		return
	}

	clients, err := ctrl.service.GetAll(c.Request.Context(), limit, offset)
	if errors.Is(err, psgrep.ErrClientNotFound) {
		ctrl.logger.Warn("Client not found", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusNotFound, gin.H{"warning": "404: client not found"})
		return
	}

	if err != nil {
		ctrl.logger.Error("Failed to retrieve clients", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Failed to retrieve client"})
		return
	}

	clientDTOs := make([]dto.ClientDTO, len(clients), cap(clients))

	for i, client := range clients {
		dto, err := mapper.ClientToDTO(&client)
		if err != nil {
			ctrl.logger.Error("Mapping error! Something is waste!", logger.Err(err), "op", op)
			continue
		}

		clientDTOs[i] = dto
	}

	ctrl.logger.Debug("Retrieved all clients", "limit", limit, "offset", offset, "op", op)
	ctrl.responce(c, http.StatusOK, clientDTOs)
}

// Get /api/v1/clients?name=&surname=
// Search client by name and surname
// 200
// 400
// 404
func (ctrl *ClientController) GetByNameAndSurname(c *gin.Context) {
	op := "controllers.clientController.getByNameAndSurname"
	name := c.Query("name")
	surname := c.Query("surname")

	clients, err := ctrl.service.GetByNameAndSurname(c.Request.Context(), name, surname)
	if errors.Is(err, psgrep.ErrClientNotFound) {
		ctrl.logger.Warn("Client not found", "op", op)
		ctrl.responce(c, http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	if err != nil {
		ctrl.logger.Error("Error while searching client ", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Name and surname cannot be empty!"})
		return
	}

	clientDTO := make([]dto.ClientDTO, len(clients), cap(clients))

	for i, client := range clients {
		dto, err := mapper.ClientToDTO(&client)
		if err != nil {
			ctrl.logger.Error("Mapping error! Somthing is waste!", logger.Err(err), "op", op)
			continue
		}

		clientDTO[i] = dto
	}

	ctrl.logger.Debug("Client retrieved successfully", "name", name, "surname", surname, "op", op)
	ctrl.responce(c, http.StatusOK, clientDTO)
}

// Patch /api/v1/clients/:id/address
// Change a address id parameter by a given new id parameter
// 200
// 400
// 500
func (ctrl *ClientController) UpdateAddress(c *gin.Context) {
	op := "controllers.clientController.UpdateAddress"
	var updateDTO dto.UpdateAddressDTO

	if err := ctrl.mapping(c, &updateDTO); err != nil {
		ctrl.logger.Error("Failed to bind JSON for ChangeAddressIdParameter", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Invalid request payload"})
		return
	}

	client, err := mapper.UpdateAddressToClientDomain(&updateDTO)
	if err != nil {
		return
	}

	if err := ctrl.service.UpdateAddress(c.Request.Context(), &client); err != nil {
		ctrl.logger.Error("Failed to update address ID", "clientID", client.Id, logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Failed to update address ID"})
		return
	}

	ctrl.logger.Debug("Address ID updated successfully", "Client ID", client.Id, "New Address ID", client.AddressId, "op", op)
	ctrl.responce(c, http.StatusOK, client)
}

// Delete /api/v1/clients/:id
// Delete client by identificator
// 204
// 400
// 500
func (ctrl *ClientController) Delete(c *gin.Context) {
	op := "controllers.clientController.Delete"
	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		ctrl.logger.Error("Invalid id", "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	if err := ctrl.service.Delete(c.Request.Context(), uuid); err != nil {
		ctrl.logger.Error("Failed to delete client", "clientID", id, logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Failed to delete client"})
		return
	}

	ctrl.logger.Debug("Client deleted successfully", "clientID", id, "op", op)
	c.Status(http.StatusNoContent)
}
