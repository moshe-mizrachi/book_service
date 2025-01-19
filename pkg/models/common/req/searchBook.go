package req

import (
	face "book_service/pkg/interfaces"
	m "book_service/pkg/models/common"
	"book_service/pkg/utils"
	"errors"
)

var _ face.Validatable = (*SearchBooks)(nil)

func (g *SearchBooks) Validate() error {
	validRange := utils.IsValidRange(g.PriceRange)
	if !validRange {
		return errors.New("invalid priceRange")
	}
	return nil
}

type SearchBooks struct {
	Title      string       `json:"title,omitempty"`
	AuthorName string       `json:"author_name,omitempty"`
	PriceRange m.PriceRange `json:"price_range,omitempty"`
	Size       int          `form:"size" validate:"gte=0,lte=100"`
	From       int          `form:"from" validate:"gte=0,lte=9999999"`
}
