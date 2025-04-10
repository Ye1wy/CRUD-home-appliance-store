package dto

import "github.com/google/uuid"

type SupplierDTO struct {
	Name        string    `json:"name" xml:"name" binding:"required"`
	AddressId   uuid.UUID `json:"address_id" xml:"address_id" binding:"required"`
	PhoneNumber string    `json:"phone_number" xml:"phone_number" binding:"required"`
}
