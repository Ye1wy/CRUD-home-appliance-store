package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func SupplierToDTO(supplier domain.Supplier) dto.Supplier {
	address := AddressToDto(*supplier.Address)
	return dto.Supplier{
		Name:        supplier.Name,
		PhoneNumber: supplier.PhoneNumber,
		Address:     &address,
	}
}

func SupplierToDomain(supplier dto.Supplier) domain.Supplier {
	address := AddressToDomain(*supplier.Address)
	return domain.Supplier{
		Name:        supplier.Name,
		PhoneNumber: supplier.PhoneNumber,
		Address:     &address,
	}
}
