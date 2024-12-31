package common

import (
	"book_service/pkg/clients"
	"book_service/pkg/constants"
	"book_service/pkg/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ActionRoutes(router *gin.Engine) {
	actionGroup := router.Group(constants.ActionRoute)
	{
		actionGroup.GET("", func(c *gin.Context) {
			username := middlewares.GetUserName(c)
			userActions, err := clients.GetLastActions(username)
			if err != nil {
				c.JSON(http.StatusOK, userActions)
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		})
	}
}
