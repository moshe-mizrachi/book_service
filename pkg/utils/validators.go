package utils

import (
	"book_service/pkg/consts"
	m "book_service/pkg/models/common"
	"fmt"

	"github.com/samber/lo"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return lo.Ternary(err == nil, true, false)
}

func IsValidRange(price m.PriceRange) bool {
	return lo.Ternary(price.Min <= price.Max, true, false)
}

func GetValidatedPayload[T any](c *gin.Context) (T, error) {
	val, exists := c.Get(consts.ValidatedAccess)
	if !exists {
		var zero T
		return zero, fmt.Errorf("payload not found in context")
	}
	typedVal, _ := val.(T)
	return typedVal, nil
}
