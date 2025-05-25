package domain

import "github.com/google/uuid"

type Supplier struct {
	Id          uuid.UUID `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	AddressId   uuid.UUID `json:"address_id" bson:"address_id"`
	PhoneNumber string    `json:"phone_number" bson:"phone_number"`
}
