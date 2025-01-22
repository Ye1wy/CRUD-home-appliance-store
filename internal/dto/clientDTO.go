package dto

type ClientDTO struct {
	Name      string `json:"name" binding:"required"`
	Surname   string `json:"surname" binding:"required"`
	Birthday  string `json:"birthday" binding:"required"`
	Gender    string `json:"gender" binding:"required"`
	AddressID string `json:"address_id" binding:"required"`
}
