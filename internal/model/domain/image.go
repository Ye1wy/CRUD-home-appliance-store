package domain

import "github.com/google/uuid"

type Image struct {
	Id   uuid.UUID `json:"id,omitempty" bson:"_id,omitempty"`
	Hash string    `json:"hash" bson:"hash"`
	Data []byte    `json:"data" bson:"data"`
}
