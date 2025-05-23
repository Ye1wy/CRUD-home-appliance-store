package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func ProductToDTO(product domain.Product) dto.Product {
	return dto.Product{
		Name:           product.Name,
		Category:       product.Category,
		Price:          product.Price,
		AvailableStock: product.AvailableStock,
		SupplierId:     product.SupplierId,
		ImageId:        product.ImageId,
	}
}

func ProductToDomain(dto dto.Product) domain.Product {
	return domain.Product{
		Name:           dto.Name,
		Category:       dto.Category,
		Price:          dto.Price,
		AvailableStock: dto.AvailableStock,
		SupplierId:     dto.SupplierId,
		ImageId:        dto.ImageId,
	}
}
