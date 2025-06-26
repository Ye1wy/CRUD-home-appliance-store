package dto

type Image struct {
	Title string `json:"title" xml:"title" binding:"required"`
	Image []byte `json:"image" xml:"image" binding:"required"`
}
