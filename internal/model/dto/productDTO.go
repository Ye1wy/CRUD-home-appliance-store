package dto

import "github.com/google/uuid"

type ProductDTO struct {
	Name           string    `json:"name" binding:"required"`
	Category       string    `json:"category" binding:"required"`
	Price          float32   `json:"price" binding:"required"`
	AvailableStock int64     `json:"available_stock" binding:"required"`
	SupplierId     uuid.UUID `json:"supplier_id" binding:"required"`
	ImageId        uuid.UUID `json:"image_id" binding:"required"`
}
