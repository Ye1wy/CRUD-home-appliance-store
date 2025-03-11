package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
)

func ToProductDTO(product *model.Product) *dto.ProductDTO {
	return &dto.ProductDTO{
		Id:             product.Id,
		Name:           product.Name,
		Category:       product.Category,
		Price:          product.Price,
		AvailableStock: product.AvailableStock,
		SupplierId:     product.SupplierId,
		ImageId:        product.ImageId,
	}
}

func ToProductModel(dto *dto.ProductDTO) *model.Product {
	return &model.Product{
		Id:             dto.Id,
		Name:           dto.Name,
		Category:       dto.Category,
		Price:          dto.Price,
		AvailableStock: dto.AvailableStock,
		SupplierId:     dto.SupplierId,
		ImageId:        dto.ImageId,
	}
}

func ToProductDTOs(products []model.Product) []dto.ProductDTO {
	dtos := make([]dto.ProductDTO, len(products))

	for i, product := range products {
		dtos[i] = *ToProductDTO(&product)
	}

	return dtos
}

func ToProductModels(dtos []dto.ProductDTO) []model.Product {
	models := make([]model.Product, len(dtos))

	for i, dto := range dtos {
		product := ToProductModel(&dto)
		models[i] = *product
	}

	return models
}
