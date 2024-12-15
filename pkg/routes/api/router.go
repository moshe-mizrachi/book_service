package api

import (
	v1 "book_service/pkg/routes/api/v1"

	"github.com/gin-gonic/gin"
)

func ApiRouter(router *gin.Engine) {
	api := router.Group("/api")
	v1.RegisterBooksRoutes(api)
}
