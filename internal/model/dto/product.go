package dto

type Product struct {
	Name           string  `json:"name" xml:"name" binding:"required"`
	Category       string  `json:"category" xml:"category" binding:"required"`
	Price          float32 `json:"price" xml:"price" binding:"required"`
	AvailableStock int64   `json:"available_stock" xml:"available_stock" binding:"required"`
	Supplier
	Image
}
