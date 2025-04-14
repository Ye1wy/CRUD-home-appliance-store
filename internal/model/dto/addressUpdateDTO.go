package dto

import "github.com/google/uuid"

type UpdateAddressDTO struct {
	Id        uuid.UUID `json:"id" xml:"id" binding:"required"`
	AddressID uuid.UUID `json:"address_id" xml:"address_id" binding:"required"`
}

//62dfb51f-ce35-4fc2-af98-e93e510a4dbb
//ab8b3a17-0031-43e0-83c4-bca4a3ccb0cd
