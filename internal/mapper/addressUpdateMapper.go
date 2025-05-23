package mapper

import (
	crud_errors "CRUD-HOME-APPLIANCE-STORE/internal/errors"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func UpdateAddressToClientDomain(dto *dto.UpdateAddress) (domain.Client, error) {
	if dto == nil {
		return domain.Client{}, crud_errors.ErrNoContent
	}

	return domain.Client{
		AddressId: dto.AddressID,
	}, nil
}

func UpdateAddressToSupplierDomain(dto *dto.UpdateAddress) (domain.Supplier, error) {
	if dto == nil {
		return domain.Supplier{}, crud_errors.ErrNoContent
	}

	return domain.Supplier{
		AddressId: dto.AddressID,
	}, nil
}
