package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func ProductDomainToProductResponse(product domain.Product) dto.ProductResponse {
	supplier := SupplierDomainToSupplierResponse(product.Supplier)
	image := ImageDomainToImageResponse(product.Image)

	return dto.ProductResponse{
		Id:             product.Id,
		Name:           product.Name,
		Category:       product.Category,
		Price:          product.Price,
		AvailableStock: product.AvailableStock,
		Supplier:       supplier,
		Image:          image,
	}
}

func ProductRequestToDomain(request dto.ProductRequest) domain.Product {
	return domain.Product{
		Name:           request.Name,
		Category:       request.Category,
		Price:          request.Price,
		AvailableStock: request.AvailableStock,
		Supplier:       domain.Supplier{Id: request.SupplierId},
		Image:          domain.Image{Id: request.ImageId},
	}
}
