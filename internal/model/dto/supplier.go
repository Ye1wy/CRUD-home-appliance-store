package dto

type Supplier struct {
	Name        string `json:"name" xml:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" xml:"phone_number" binding:"required"`
	Address
}
