package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func SupplierToDTO(supplier *domain.Supplier) (dto.SupplierDTO, error) {
	if supplier == nil {
		return dto.SupplierDTO{}, ErrNoContent
	}

	return dto.SupplierDTO{
		Name:        supplier.Name,
		AddressId:   supplier.AddressId,
		PhoneNumber: supplier.PhoneNumber,
	}, nil
}

func SupplierToDomain(dto *dto.SupplierDTO) (domain.Supplier, error) {
	if dto == nil {
		return domain.Supplier{}, ErrNoContent
	}

	return domain.Supplier{
		Name:        dto.Name,
		AddressId:   dto.AddressId,
		PhoneNumber: dto.PhoneNumber,
	}, nil
}
