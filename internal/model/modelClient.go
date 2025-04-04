package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	Id               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ClientName       string             `json:"client_name" bson:"client_name"`
	ClientSurname    string             `json:"client_surname" bson:"client_surname"`
	Birthday         time.Time          `json:"birthday" bson:"birthday"`
	Gender           string             `json:"gender" bson:"gender"`
	RegistrationDate time.Time          `json:"registration_date" bson:"registration_date"`
	AddressId        string             `json:"address_id" bson:"address_id"`
}
