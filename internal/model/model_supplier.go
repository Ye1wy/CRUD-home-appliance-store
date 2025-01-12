package model

type Supplier struct {
	Id          int64  `json:"id,omitempty" bson:"id,omitempt"`
	Name        string `json:"name" bson:"name"`
	AddressId   int64  `json:"address_id" bson:"address_id"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}
