package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
)

func ProductToDTO(product *model.Product) (*dto.ProductDTO, error) {
	if product == nil {
		return nil, nil
	}

	return &dto.ProductDTO{
		Id:             product.Id,
		Name:           product.Name,
		Category:       product.Category,
		Price:          product.Price,
		AvailableStock: product.AvailableStock,
		SupplierId:     product.SupplierId,
		ImageId:        product.ImageId,
	}, nil
}

func ProductToModel(dto *dto.ProductDTO) (*model.Product, error) {
	if dto == nil {
		return nil, nil
	}

	return &model.Product{
		Id:             dto.Id,
		Name:           dto.Name,
		Category:       dto.Category,
		Price:          dto.Price,
		AvailableStock: dto.AvailableStock,
		SupplierId:     dto.SupplierId,
		ImageId:        dto.ImageId,
	}, nil
}
