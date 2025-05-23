package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func SupplierToDTO(supplier domain.Supplier) dto.Supplier {
	return dto.Supplier{
		Name:        supplier.Name,
		AddressId:   supplier.AddressId,
		PhoneNumber: supplier.PhoneNumber,
	}
}

func SupplierToDomain(dto dto.Supplier) domain.Supplier {
	return domain.Supplier{
		Name:        dto.Name,
		AddressId:   dto.AddressId,
		PhoneNumber: dto.PhoneNumber,
	}
}
