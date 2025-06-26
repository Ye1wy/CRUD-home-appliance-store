package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func ImageToDomain(dto dto.Image) domain.Image {
	return domain.Image{
		Title: dto.Title,
		Data:  dto.Image,
	}
}

func ImageToDTO(domain domain.Image) dto.Image {
	return dto.Image{
		Title: domain.Title,
		Image: domain.Data,
	}
}
