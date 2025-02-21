package dto

type ProductDTO struct {
	Name           string  `json:"name" binding:"required"`
	Category       string  `json:"category" binding:"required"`
	Price          float32 `json:"price" binding:"required"`
	AvailableStock int64   `json:"available_stock" binding:"required"`
	SupplierId     string  `json:"supplier_id" binding:"required"`
	ImageId        string  `json:"image_id" binding:"required"`
}
