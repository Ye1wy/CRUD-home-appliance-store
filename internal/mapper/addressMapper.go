package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func AddressToDomain(dto dto.Address) domain.Address {
	return domain.Address{
		Country: dto.Country,
		City:    dto.City,
		Street:  dto.Street,
	}
}

func AddressToDto(domain domain.Address) dto.Address {
	return dto.Address{
		Country: domain.Country,
		City:    domain.City,
		Street:  domain.Street,
	}
}
