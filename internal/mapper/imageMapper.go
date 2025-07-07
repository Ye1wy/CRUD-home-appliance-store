package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
)

func ImageRequestToDomain(dto dto.ImageRequest) domain.Image {
	return domain.Image{
		Title: dto.Title,
		Data:  dto.Image,
	}
}

func ImageDomainToImageResponse(domain domain.Image) dto.ImageResponse {
	return dto.ImageResponse{
		Id:    domain.Id,
		Title: domain.Title,
		Image: domain.Data,
	}
}
