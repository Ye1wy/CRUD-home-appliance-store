package dto

import "github.com/google/uuid"

type ProductRequest struct {
	Name           string    `json:"name" xml:"name" binding:"required"`
	Category       string    `json:"category" xml:"category" binding:"required"`
	Price          float32   `json:"price" xml:"price" binding:"required"`
	AvailableStock int64     `json:"available_stock" xml:"available_stock" binding:"required"`
	SupplierId     uuid.UUID `json:"supplier_id" xml:"supplier_id" binding:"required"`
	ImageId        uuid.UUID `json:"image_id" xml:"image_id" binding:"required"`
}

type ProductResponse struct {
	Id             uuid.UUID        `json:"id" xml:"id"`
	Name           string           `json:"name" xml:"name"`
	Category       string           `json:"category" xml:"category"`
	Price          float32          `json:"price" xml:"price"`
	AvailableStock int64            `json:"available_stock" xml:"available_stock"`
	Supplier       SupplierResponse `json:"supplier" xml:"supplier"`
	Image          ImageResponse    `json:"image" xml:"image"`
}
