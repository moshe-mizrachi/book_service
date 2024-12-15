package common

import "github.com/google/uuid"

type Book struct {
	ID             uuid.UUID `json:"-"`
	Title          string    `json:"title" validate:"required,min=2,max=250"`
	AuthorName     string    `json:"author_name" validate:"required,min=2,max=100"`
	Price          float64   `json:"price" validate:"required,gte=0,lte=10000"`
	EbookAvailable bool      `json:"ebook_available" validate:"required"`
	PublishDate    string    `json:"publish_date" validate:"required" copier:"-"`
}

type PriceRange struct {
	Min float64 `json:"min" validate:"gte=0,lte=10000"`
	Max float64 `json:"max" validate:"gte=0,lte=10000"`
}

type TitleUpdate struct {
	Title string `json:"title"`
}
