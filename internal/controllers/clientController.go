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
	Create(ctx context.Context, client *domain.Client) error
	GetAll(ctx context.Context, limit, offset int) ([]domain.Client, error)
	GetByNameAndSurname(ctx context.Context, name, surname string) ([]domain.Client, error)
	UpdateAddress(ctx context.Context, id uuid.UUID, address *domain.Address) error
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
	var input dto.Client

	if err := c.ShouldBind(&input); err != nil {
		ctrl.logger.Warn("Failed to bind JSON/XML for create", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: invalid data received"})
		return
	}

	ctrl.logger.Debug("check data", "data", input)

	client, err := mapper.ClientToDomain(input)
	if err != nil {
		ctrl.logger.Warn("Failed mapping dto to domain", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid birthady date in request payload"})
		return
	}

	if err := ctrl.service.Create(c.Request.Context(), &client); err != nil {
		ctrl.logger.Error("Failed to add client", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	ctrl.logger.Debug("Client created", "op", op)
	c.Status(http.StatusCreated)
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

	clients, err := ctrl.service.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		if errors.Is(err, crud_errors.ErrInvalidParam) {
			ctrl.logger.Warn("Invalid limit or offset parameter", "op", op)
			ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: limit cannot be less or equal 0, offset cannot be less than 0"})
			return
		}

		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("No content", logger.Err(err), "op", op)
			ctrl.responce(c, http.StatusNotFound, gin.H{"massage": "404: no data is contains"})
			return
		}

		ctrl.logger.Error("Failed to retrieve clients", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	output := make([]dto.Client, len(clients))

	for i, client := range clients {
		dto := mapper.ClientToDTO(client)
		output[i] = dto
	}

	ctrl.logger.Debug("Retrieved all client's", "limit", limit, "offset", offset, "op", op)
	ctrl.responce(c, http.StatusOK, output)
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
	if err != nil {
		if errors.Is(err, crud_errors.ErrInvalidParam) {
			ctrl.logger.Warn("Invalid value name or surname", logger.Err(err), "op", op)
			ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalid request payload: name and surname cannot be empty"})
			return
		}

		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("Client not found", "op", op)
			ctrl.responce(c, http.StatusNotFound, gin.H{"massage": "404: client not found"})
			return
		}

		ctrl.logger.Error("Failed to get data from database", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	output := make([]dto.Client, len(clients))

	for i, client := range clients {
		dto := mapper.ClientToDTO(client)
		output[i] = dto
	}

	ctrl.logger.Debug("Client retrieved", "name", name, "surname", surname, "op", op)
	ctrl.responce(c, http.StatusOK, output)
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
//	@Router			/api/v1/clients/:id [patch]
func (ctrl *ClientController) UpdateAddress(c *gin.Context) {
	op := "controllers.clientController.UpdateAddress"
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

	newAddress := mapper.AddressToDomain(input)

	if err := ctrl.service.UpdateAddress(c.Request.Context(), id, &newAddress); err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("Client not found", logger.Err(err), "op", op)
			ctrl.responce(c, http.StatusNotFound, gin.H{"massage": "404: client not found for update"})
			return
		}

		ctrl.logger.Error("Failed to update address ID", "id", id, logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	ctrl.logger.Debug("Client updated", "id", id, "op", op)
	c.Status(http.StatusOK)
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
		ctrl.logger.Warn("The received identifier is invalid", "op", op)
		ctrl.responce(c, http.StatusBadRequest, gin.H{"massage": "Invalud request payload: id is not valid"})
		return
	}

	if err := ctrl.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			ctrl.logger.Debug("Client not found", "op", op)
			c.Status(http.StatusNoContent)
			return
		}

		ctrl.logger.Error("Failed delete client by id", logger.Err(err), "op", op)
		ctrl.responce(c, http.StatusInternalServerError, gin.H{"error": "Server is busy"})
		return
	}

	ctrl.logger.Debug("Client deleted", "id", id, "op", op)
	c.Status(http.StatusNoContent)
}
