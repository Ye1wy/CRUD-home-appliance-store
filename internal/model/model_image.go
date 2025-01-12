package model

type Image struct {
	Id    string `json:"id,omitempty" bson:"id,omitempty"`
	Image []byte `json:"image" bson:"image"`
}

type ProductImage struct {
	Image []byte `json:"image" bson:"image"`
}
