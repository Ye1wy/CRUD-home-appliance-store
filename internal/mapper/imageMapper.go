package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func ImageToDomain(dto dto.Image) domain.Image {
	return domain.Image{
		Data: dto.Image,
	}
}

func ImageToDTO(domain domain.Image) dto.Image {
	return dto.Image{
		Image: domain.Data,
	}
}
