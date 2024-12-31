package req

import (
	face "book_service/pkg/interfaces"
	m "book_service/pkg/models/common"
	"book_service/pkg/utils"
	"errors"
	"time"
)

var _ face.Validatable = (*GetBook)(nil)

func (g *GetBook) Validate() error {
	validUUID := utils.IsValidUUID(g.ID)
	if !validUUID {
		return errors.New("invalid uuid")
	}
	return nil
}

func (g *SearchBooks) Validate() error {
	validRange := utils.IsValidRange(g.PriceRange)
	if !validRange {
		return errors.New("invalid priceRange")
	}
	return nil
}

type AddBook struct {
	Title          string    `json:"title" validate:"required,min=2,max=250"`
	AuthorName     string    `json:"author_name" validate:"required,min=2,max=40"`
	Price          float64   `json:"price" validate:"required,gte=0,lte=10000"`
	EbookAvailable *bool     `json:"ebook_available" validate:"required"`
	PublishDate    time.Time `json:"publish_date" validate:"required"`
}

type UpdateBook struct {
	ID    string `uri:"id" binding:"required" validate:"required"`
	Title string `json:"title" validate:"required,min=1"`
}

type GetBook struct {
	ID string `uri:"id" binding:"required"`
}

type DeleteBook struct {
	ID string `uri:"id" binding:"required" validate:"required"`
}

type SearchBooks struct {
	Title      string       `json:"title,omitempty"`
	AuthorName string       `json:"author_name,omitempty"`
	PriceRange m.PriceRange `json:"price_range,omitempty"`
	Size       int          `form:"size" validate:"gte=0,lte=100"`
	From       int          `form:"from" validate:"gte=0,lte=9999999"`
}
