package req

import (
	"time"
)

type AddBook struct {
	Title          string    `json:"title" validate:"required,min=2,max=150"`
	AuthorName     string    `json:"author_name" validate:"required,min=2,max=40"`
	Price          float64   `json:"price" validate:"required,gte=0,lte=10000"`
	EbookAvailable bool      `json:"ebook_available" validate:"required"`
	PublishDate    time.Time `json:"publish_date" validate:"required"`
	Username       string    `json:"username" validate:"omitempty,min=3,max=50"`
}

type UpdateBook struct {
	ID       string `json:"id" validate:"required"`
	Title    string `json:"title" validate:"required,min=1"`
	Username string `json:"username"`
}

type GetBook struct {
	ID       string `uri:"id" binding:"required"`
	Username string `json:"username"`
}

type DeleteBook struct {
	ID       string `json:"id" validate:"required"`
	Username string `json:"username"`
}

type SearchBooks struct {
	Title      string  `json:"title,omitempty"`
	AuthorName string  `json:"author_name,omitempty"`
	PriceRange *string `json:"price_range,omitempty"` // e.g., "10-20"
	Username   string  `json:"username"`
}

type StoreSummary struct {
	Username string `json:"username"`
}
