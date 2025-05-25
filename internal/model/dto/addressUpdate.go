package dto

import "github.com/google/uuid"

type UpdateAddress struct {
	AddressID uuid.UUID `json:"address_id" xml:"address_id" binding:"required"`
}
