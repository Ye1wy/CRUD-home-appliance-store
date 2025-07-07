package dto

import "github.com/google/uuid"

type SupplierRequest struct {
	Name        string `json:"name" xml:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" xml:"phone_number" binding:"required"`
	*Address
}

type SupplierResponse struct {
	Id          uuid.UUID `json:"id" xml:"id"`
	Name        string    `json:"name" xml:"name"`
	PhoneNumber string    `json:"phone_number" xml:"phone_number"`
	*Address
}
