package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func ProductToDTO(product domain.Product) dto.Product {
	supplier := SupplierToDTO(product.Supplier)
	image := ImageToDTO(product.Image)

	return dto.Product{
		Name:           product.Name,
		Category:       product.Category,
		Price:          product.Price,
		AvailableStock: product.AvailableStock,
		Supplier:       supplier,
		Image:          image,
	}
}

func ProductToDomain(dto dto.Product) domain.Product {
	supplier := SupplierToDomain(dto.Supplier)
	image := ImageToDomain(dto.Image)

	return domain.Product{
		Name:           dto.Name,
		Category:       dto.Category,
		Price:          dto.Price,
		AvailableStock: dto.AvailableStock,
		Supplier:       supplier,
		Image:          image,
	}
}
