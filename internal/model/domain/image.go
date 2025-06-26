package domain

import "github.com/google/uuid"

type Image struct {
	Id    uuid.UUID `json:"id,omitempty" bson:"_id,omitempty"`
	Title string    `json:"title" bson:"title"`
	Data  []byte    `json:"data" bson:"data"`
	// Hash string    `json:"hash" bson:"hash"`
}
