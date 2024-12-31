package common

import (
	"github.com/gin-gonic/gin"
)

func HealthRoutes(router *gin.Engine) {
	healthGroup := router.Group("/ping")
	{
		healthGroup.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Healthy!"})
		})
	}
}
