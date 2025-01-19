package req

import "time"

type AddBook struct {
	Title          string    `json:"title" validate:"required,min=2,max=250"`
	AuthorName     string    `json:"author_name" validate:"required,min=2,max=40"`
	Price          float64   `json:"price" validate:"required,gte=0,lte=10000"`
	EbookAvailable *bool     `json:"ebook_available" validate:"required"`
	PublishDate    time.Time `json:"publish_date" validate:"required"`
}
