package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func SupplierToDTO(supplier domain.Supplier) dto.Supplier {
	return dto.Supplier{
		Name: supplier.Name,
		// AddressId:   supplier.AddressId,
		PhoneNumber: supplier.PhoneNumber,
		Address: dto.Address{
			Country: supplier.Address.Country,
			City:    supplier.Address.City,
			Street:  supplier.Address.Street,
		},
	}
}

func SupplierToDomain(supplier dto.Supplier) domain.Supplier {
	return domain.Supplier{
		Name: supplier.Name,
		// AddressId:   dto.AddressId,
		PhoneNumber: supplier.PhoneNumber,
		Address: domain.Address{
			Country: supplier.Address.Country,
			City:    supplier.Address.City,
			Street:  supplier.Address.Street,
		},
	}
}
