package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func UpdateAddressToClientDomain(dto *dto.UpdateAddress) (domain.Client, error) {
	if dto == nil {
		return domain.Client{}, ErrNoContent
	}

	return domain.Client{
		AddressId: dto.AddressID,
	}, nil
}

func UpdateAddressToSupplierDomain(dto *dto.UpdateAddress) (domain.Supplier, error) {
	if dto == nil {
		return domain.Supplier{}, ErrNoContent
	}

	return domain.Supplier{
		AddressId: dto.AddressID,
	}, nil
}
