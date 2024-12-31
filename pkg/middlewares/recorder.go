package middlewares

import (
	"book_service/pkg/clients"
	"book_service/pkg/constants"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func RecordActions() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isActionRoute(c.Request.URL.Path) {
			c.Next()
			return
		}
		action := clients.UserAction{
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Time:   time.Now(),
		}
		username := GetUserName(c)
		action.User = username

		clients.AppendAction(action)

		c.Next()
	}
}

func GetUserName(c *gin.Context) string {
	username := c.GetHeader("X-Username")

	if username == "" {
		username = c.ClientIP()
	}
	return username
}

func isActionRoute(path string) bool {
	segments := strings.Split(path, "/")
	lastIndex := len(segments) - 1
		if segments[lastIndex] == constants.ActionRoute {
			return true
	}
	return false
}

