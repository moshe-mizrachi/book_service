package common

import (
	"book_service/pkg/clients"
	"book_service/pkg/consts"
	"book_service/pkg/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ActionRoutes(router *gin.Engine) {
	actionGroup := router.Group("/" + consts.ActionRoute)
	{
		actionGroup.GET("", func(c *gin.Context) {
			username := middlewares.GetUserName(c)
			userActions, err := clients.GetLastActions(username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, userActions)
		})
	}
}
