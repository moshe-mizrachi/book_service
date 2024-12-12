package common

import (
	"book_service/pkg/clients"
	"book_service/pkg/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ActionRoutes(router *gin.Engine) {
	actionGroup := router.Group("/action")
	{
		actionGroup.GET("/", func(c *gin.Context) {
			username := middlewares.GetUserName(c)
			userActions := clients.GetLastActions(username)
			c.JSON(http.StatusOK, userActions)
		})
	}
}