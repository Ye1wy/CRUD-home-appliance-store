package domain

import "github.com/google/uuid"

type Supplier struct {
	Id          uuid.UUID `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	PhoneNumber string    `json:"phone_number" bson:"phone_number"`
	Address     Address   `json:"address" bson:"address"`
}
