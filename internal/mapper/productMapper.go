package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func ProductToDTO(product *domain.Product) (dto.ProductDTO, error) {
	if product == nil {
		return dto.ProductDTO{}, ErrNoContent
	}

	return dto.ProductDTO{
		Name:           product.Name,
		Category:       product.Category,
		Price:          product.Price,
		AvailableStock: product.AvailableStock,
		SupplierId:     product.SupplierId,
		ImageId:        product.ImageId,
	}, nil
}

func ProductToDomain(dto *dto.ProductDTO) (domain.Product, error) {
	if dto == nil {
		return domain.Product{}, ErrNoContent
	}

	return domain.Product{
		Name:           dto.Name,
		Category:       dto.Category,
		Price:          dto.Price,
		AvailableStock: dto.AvailableStock,
		SupplierId:     dto.SupplierId,
		ImageId:        dto.ImageId,
	}, nil
}
