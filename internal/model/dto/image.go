package dto

import "github.com/google/uuid"

type ImageRequest struct {
	Title string `json:"title" xml:"title" binding:"required"`
	Image []byte `json:"image" xml:"image" binding:"required"`
}

type ImageResponse struct {
	Id    uuid.UUID `json:"id" xml:"id"`
	Title string    `json:"title" xml:"title"`
	Image []byte    `json:"image" xml:"image"`
}
