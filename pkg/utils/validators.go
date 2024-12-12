package utils

import (
	m "book_service/pkg/models/common"

	"github.com/google/uuid"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func IsValidRange(price m.PriceRange) bool {
	return price.Min > price.Max
}
