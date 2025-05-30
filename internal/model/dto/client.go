package dto

type Client struct {
	Name     string `json:"name" xml:"name" binding:"required"`
	Surname  string `json:"surname" xml:"surname" binding:"required"`
	Birthday string `json:"birthday" xml:"birthday" binding:"required"`
	Gender   string `json:"gender" xml:"gender" binding:"required"`
	Address
}
