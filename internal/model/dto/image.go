package dto

type Image struct {
	Image []byte `json:"image" xml:"image" binding:"required"`
}
