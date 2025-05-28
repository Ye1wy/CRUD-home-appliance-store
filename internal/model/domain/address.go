package domain

import "github.com/google/uuid"

type Address struct {
	Id      uuid.UUID `json:"id,omitempty" bson:"_id,omitempty"`
	Country string    `json:"country" bson:"country"`
	City    string    `json:"city" bson:"city"`
	Street  string    `json:"street" bson:"street"`
}
