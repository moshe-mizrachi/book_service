package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

var (
	validate    = validator.New()
	binderFuncs = []func(*gin.Context, any) error{
		func(c *gin.Context, obj any) error { return c.ShouldBindUri(obj) },
		func(c *gin.Context, obj any) error { return c.ShouldBindQuery(obj) },
		func(c *gin.Context, obj any) error { return c.ShouldBindJSON(obj) },
	}
)

func canIgnoreError(err error) bool {
	return errors.Is(err, io.EOF)
}

func Validation[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload T

		for _, binder := range binderFuncs {
			if err := binder(c, &payload); err != nil && !canIgnoreError(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
		}

		if err := validate.Struct(payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("validated", payload)
		c.Next()
	}
}
