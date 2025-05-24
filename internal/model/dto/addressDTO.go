package dto

type Address struct {
	Country string `json:"country" xml:"country"`
	City    string `json:"city" xml:"city"`
	Street  string `json:"street" xml:"street"`
}
