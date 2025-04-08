package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Supplier struct {
	Id          primitive.ObjectID `json:"_id"`
	Name        string             `json:"name"`
	AddressId   string             `json:"address_id"`
	PhoneNumber string             `json:"phone_number"`
}
