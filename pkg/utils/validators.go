package utils

import (
	"book_service/pkg/constants"
	m "book_service/pkg/models/common"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func IsValidRange(price m.PriceRange) bool {
	return price.Min <= price.Max
}

func GetValidatedPayload[T any](c *gin.Context) T {
	val, exists := c.Get(constants.ValidatedAccess)
	if !exists {
		var zero T
		return zero
	}
	typedVal, _ := val.(T)
	return typedVal
}
