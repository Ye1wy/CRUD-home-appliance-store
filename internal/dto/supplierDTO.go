package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type SupplierDTO struct {
	Id          primitive.ObjectID `json:"id" xml:"id"`
	Name        string             `json:"name" xml:"name" binding:"required"`
	AddressId   string             `json:"address_id" xml:"address_id" binding:"required"`
	PhoneNumber string             `json:"phone_number" xml:"phone_number" binding:"required"`
}
