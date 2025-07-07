package dto

import "github.com/google/uuid"

type ClientRequest struct {
	Name     string `json:"name" xml:"name" binding:"required"`
	Surname  string `json:"surname" xml:"surname" binding:"required"`
	Birthday string `json:"birthday" xml:"birthday" binding:"required"`
	Gender   string `json:"gender" xml:"gender" binding:"required"`
	*Address
}

type ClientResponse struct {
	Id       uuid.UUID `json:"id" xml:"id"`
	Name     string    `json:"name" xml:"name"`
	Surname  string    `json:"surname" xml:"surname"`
	Birthday string    `json:"birthday" xml:"birthday"`
	Gender   string    `json:"gender" xml:"gender"`
	*Address
}
