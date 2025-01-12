package model

type Client struct {
	Id               int64  `json:"id,omitempty" bson:"id,omitempty"`
	ClientName       string `json:"client_name" bson:"client_name"`
	ClientSurname    string `json:"client_surname" bson:"client_surname"`
	Birthday         string `json:"birthday" bson:"birthday"`
	Gender           string `json:"gender" bson:"gender"`
	RegistrationDate string `json:"registration_date" bson:"registration_date"`
	AddressId        int64  `json:"address_id" bson:"address_id"`
}
