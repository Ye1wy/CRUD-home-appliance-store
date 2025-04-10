package dto

import "github.com/google/uuid"

type UpdateAddressDTO struct {
	Id        uuid.UUID `json:"id" xml:"id" binding:"required"`
	AddressID uuid.UUID `json:"address_id" xml:"address_id" binding:"required"`
}
