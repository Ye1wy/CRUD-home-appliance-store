package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id             uuid.UUID `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string    `json:"name" bson:"name"`
	Category       string    `json:"category" bson:"category"`
	Price          float32   `json:"price" bson:"price"`
	AvailableStock int64     `json:"available_stock" bson:"available_stock"`
	LastUpdateDate time.Time `json:"last_update_date" bson:"last_update_date"`
	Supplier       Supplier  `json:"supplier" bson:"supplier"`
	Image          Image     `json:"image" bson:"image"`
}
