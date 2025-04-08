package dto

type UpdateAddressDTO struct {
	AddressID string `json:"address_id" xml:"address_id" binding:"required"`
}
