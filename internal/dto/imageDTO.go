package dto

type ImageDTO struct {
	Image []byte `json:"image" bson:"image"`
}
