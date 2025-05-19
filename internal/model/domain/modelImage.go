package domain

import "github.com/google/uuid"

type Image struct {
	Id    uuid.UUID `json:"id,omitempty" bson:"_id,omitempty"`
	Image []byte    `json:"image"`
}
