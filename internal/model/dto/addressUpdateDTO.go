package dto

import "github.com/google/uuid"

type UpdateAddressDTO struct {
	AddressID uuid.UUID `json:"address_id" xml:"address_id" binding:"required"`
}
