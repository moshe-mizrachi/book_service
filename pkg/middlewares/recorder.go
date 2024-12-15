package middlewares

import (
	"book_service/pkg/clients"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func isActionRoute(path string) bool {
	segments := strings.Split(path, "/")

	for _, segment := range segments {
		if segment == "action" {
			return true
		}
	}
	return false
}

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
