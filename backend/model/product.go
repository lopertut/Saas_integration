package model

type Product struct {
	Id             int            `json:"id"`
	Name           string         `json:"name"`
	Price          float64        `json:"price"`
	Specifications map[string]any `json:"specs"`
	Images         []string       `json:"images"`
	CategoryName   string         `json:"category"`
	BrandName      string         `json:"brand"`
	AvgRating      float64        `json:"rating"`
}
