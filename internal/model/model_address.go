package model

type Address struct {
	ID      int64  `json:"id,omitempty" bson:"id,omitempty"`
	Country string `json:"country" bson:"country"`
	City    string `json:"city" bson:"city"`
	Street  string `json:"street" bson:"street"`
}
