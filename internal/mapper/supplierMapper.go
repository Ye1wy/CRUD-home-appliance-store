package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func SupplierDomainToSupplierResponse(supplier domain.Supplier) dto.SupplierResponse {
	address := AddressToDto(*supplier.Address)
	return dto.SupplierResponse{
		Id:          supplier.Id,
		Name:        supplier.Name,
		PhoneNumber: supplier.PhoneNumber,
		Address:     &address,
	}
}

func SupplierRequestToDomain(supplier dto.SupplierRequest) domain.Supplier {
	address := AddressToDomain(*supplier.Address)
	return domain.Supplier{
		Name:        supplier.Name,
		PhoneNumber: supplier.PhoneNumber,
		Address:     &address,
	}
}
