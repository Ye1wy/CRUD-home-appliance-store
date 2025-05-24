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

type clientService interface {
	Create(ctx context.Context, client domain.Client, address domain.Address) error
	GetAll(ctx context.Context, limit, offset int) ([]domain.Client, error)
	GetByNameAndSurname(ctx context.Context, name, surname string) ([]domain.Client, error)
	UpdateAddress(ctx context.Context, id uuid.UUID, address domain.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ClientController struct {
	*BaseController
	service clientService
}

func NewClientsController(service clientService, logger *logger.Logger) *ClientController {
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
//		@Failure		400	{object}	domain.Error
//		@Failure		404	{object}	domain.Error
//		@Failure		500	{object}	domain.Error
//		@Router			/api/v1/clients [post]
func (ctrl *ClientController) Create(c *gin.Context) {
	op := "controllers.clientController.Create"
	var inputData dto.FullClientData

	if err := c.ShouldBind(&inputData); err != nil {
		ctrl.logger.Error("Failed to bind JSON for create", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	client, err := mapper.ClientToDomain(inputData.Client)
	if err != nil {
		ctrl.logger.Error("Failed mapping dto to domain", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Invalid birthady date in request payload"})
		return
	}

	address := mapper.AddressToDomain(inputData.Address)

	// var clientDTO dto.Client

	// if err := c.ShouldBind(&clientDTO); err != nil {
	// 	ctrl.logger.Error("Failed to bind JSON for Create", logger.Err(err), "op", op)
	// 	ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
	// 	return
	// }

	// client, err := mapper.ClientToDomain(clientDTO)
	// if err != nil {
	// 	ctrl.logger.Error("Failed mapping dto to domain", logger.Err(err), "op", op)
	// 	ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Error from server"})
	// 	return
	// }

	if err := ctrl.service.Create(c.Request.Context(), client, address); err != nil {
		ctrl.logger.Error("Failed to add client: ", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Failed to add client"})
		return
	}

	ctrl.logger.Info("Client added successfully", "op", op)
	ctrl.responce(c, http.StatusCreated, client)
}

// Get All Client godoc
//
//	@Summary		Get all client
//	@Description	That methods retrive all registered client in system
//	@Tags			clients
//	@Accept			json/xml
//	@Produce		json/xml
//	@Success		200	{object}	dto.ClientDTO
//	@Failure		400	{object}	domain.Error
//	@Failure		404	{object}	domain.Error
//	@Failure		500	{object}	domain.Error
//	@Router			/api/v1/clients [get]
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
	if errors.Is(err, crud_errors.ErrNotFound) {
		ctrl.logger.Warn("Client not found", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusNotFound, gin.H{"warning": "404: client not found"})
		return
	}

	if err != nil {
		ctrl.logger.Error("Failed to retrieve clients", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Failed to retrieve client"})
		return
	}

	ctrl.logger.Debug("aboba", "clients", clients)
	clientDTOs := make([]dto.Client, len(clients), cap(clients))

	for i, client := range clients {
		dto := mapper.ClientToDTO(client)
		clientDTOs[i] = dto
	}

	ctrl.logger.Debug("Retrieved all clients", "limit", limit, "offset", offset, "op", op)
	ctrl.responce(c, http.StatusOK, clientDTOs)
}

// Get Client godoc
//
//	@Summary		Get client filtered by name and surname
//	@Description	That methods retrive all required registered client in system
//	@Tags			clients
//	@Accept			json/xml
//	@Produce		json/xml
//	@Success		200	{object}	[]dto.ClientDTO
//	@Failure		400	{object}	domain.Error
//	@Failure		404	{object}	domain.Error
//	@Failure		500	{object}	domain.Error
//	@Router			/api/v1/clients/search [get]
func (ctrl *ClientController) GetByNameAndSurname(c *gin.Context) {
	op := "controllers.clientController.getByNameAndSurname"
	name := c.Query("name")
	surname := c.Query("surname")

	clients, err := ctrl.service.GetByNameAndSurname(c.Request.Context(), name, surname)
	if errors.Is(err, crud_errors.ErrNotFound) {
		ctrl.logger.Warn("Client not found", "op", op)
		ctrl.responce(c, http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	if err != nil {
		ctrl.logger.Error("Error while searching client ", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Name and surname cannot be empty!"})
		return
	}

	clientDTO := make([]dto.Client, len(clients), cap(clients))

	for i, client := range clients {
		dto := mapper.ClientToDTO(client)
		clientDTO[i] = dto
	}

	ctrl.logger.Debug("Client retrieved successfully", "name", name, "surname", surname, "op", op)
	ctrl.responce(c, http.StatusOK, clientDTO)
}

// Update Client field godoc
//
//	@Summary		Update address on client
//	@Description	That methods change address on client
//	@Tags			clients
//	@Accept			json/xml
//	@Produce		json/xml
//	@Param	address_id path string true "New address"
//	@Success		200	{object}
//	@Failure		400	{object}	domain.Error
//	@Failure		404	{object}	domain.Error
//	@Failure		500	{object}	domain.Error
//	@Router			/api/v1/clients/:id/decrease [patch]
func (ctrl *ClientController) UpdateAddress(c *gin.Context) {
	op := "controllers.clientController.UpdateAddress"
	rawId := c.Param("id")
	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Error("Failed parse id to uuid", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Invalud request payload"})
		return
	}

	var updateDTO dto.Address

	if err := c.ShouldBind(&updateDTO); err != nil {
		ctrl.logger.Error("Failed to bind JSON for ChangeAddressIdParameter", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Invalid request payload"})
		return
	}

	recentData := mapper.AddressToDomain(updateDTO)

	if err := ctrl.service.UpdateAddress(c.Request.Context(), id, recentData); err != nil {
		ctrl.logger.Error("Failed to update address ID", "clientID", id, logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Failed to update address ID"})
		return
	}

	ctrl.logger.Debug("Address ID updated successfully", "op", op)
	ctrl.responce(c, http.StatusOK, gin.H{"massage": "address updated"})
}

// Delete Client godoc
//
//	@Summary		Delete client from system
//	@Description	That methods deleting registered client in system by id
//	@Tags			clients
//	@Accept			json/xml
//	@Produce		json/xml
//	@Success		204	{object}
//	@Failure		400	{object}	domain.Error
//	@Failure		500	{object}	domain.Error
//	@Router			/api/v1/clients/:id [delete]
func (ctrl *ClientController) Delete(c *gin.Context) {
	op := "controllers.clientController.Delete"
	rawId := c.Param("id")

	id, err := uuid.Parse(rawId)
	if err != nil {
		ctrl.logger.Error("Invalid id", "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	if err := ctrl.service.Delete(c.Request.Context(), id); err != nil {
		ctrl.logger.Error("Failed to delete client", "clientID", rawId, logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Failed to delete client"})
		return
	}

	ctrl.logger.Debug("Client deleted successfully", "clientID", rawId, "op", op)
	c.Status(http.StatusNoContent)
}
