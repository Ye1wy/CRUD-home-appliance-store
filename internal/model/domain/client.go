package domain

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	Id               uuid.UUID `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string    `json:"client_name" bson:"client_name"`
	Surname          string    `json:"client_surname" bson:"client_surname"`
	Birthday         time.Time `json:"birthday" bson:"birthday"`
	Gender           string    `json:"gender" bson:"gender"`
	RegistrationDate time.Time `json:"registration_date" bson:"registration_date"`
	Address          *Address  `json:"address,omitempty" bson:"address,omitempty"`
}
