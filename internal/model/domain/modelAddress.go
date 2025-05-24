package domain

import "github.com/google/uuid"

type Address struct {
	Id      uuid.UUID `json:"id,omitempty"`
	Country string    `json:"country"`
	City    string    `json:"city"`
	Street  string    `json:"street"`
}
