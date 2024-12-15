package common

import (
	v1 "book_service/pkg/handlers/v1"
	"github.com/gin-gonic/gin"
)

func StatisticRoutes(router *gin.Engine) {
	statGroup := router.Group("/store")
	{
		statGroup.GET("/", v1.GetBooksStats)
	}
}
