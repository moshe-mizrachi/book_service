package common

import "github.com/google/uuid"

type Book struct {
	ID             uuid.UUID `json:"-"`
	Title          string    `json:"title" validate:"required,min=2,max=150"`
	AuthorName     string    `json:"author_name" validate:"required,min=2,max=40"`
	Price          float64   `json:"price" validate:"required,gte=0,lte=10000"`
	EbookAvailable bool      `json:"ebook_available" validate:"required"`
	PublishDate    string    `json:"publish_date" validate:"required" copier:"-"`
	Username       string    `json:"username" validate:"omitempty,min=3,max=50"`
}
