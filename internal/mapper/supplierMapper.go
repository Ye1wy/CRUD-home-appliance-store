package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
)

func SupplierToDTO(supplier *model.Supplier) (*dto.SupplierDTO, error) {
	if supplier == nil {
		return nil, nil
	}

	return &dto.SupplierDTO{
		Id:          supplier.Id,
		Name:        supplier.Name,
		AddressId:   supplier.AddressId,
		PhoneNumber: supplier.PhoneNumber,
	}, nil
}

func SupplierToModel(dto *dto.SupplierDTO) (*model.Supplier, error) {
	if dto == nil {
		return nil, nil
	}

	return &model.Supplier{
		Id:          dto.Id,
		Name:        dto.Name,
		AddressId:   dto.AddressId,
		PhoneNumber: dto.PhoneNumber,
	}, nil
}
