package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process the request
		c.Next()

		// Log request details
		logrus.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"latency":    time.Since(start),
			"user_agent": c.Request.UserAgent(),
			"error":      c.Errors.ByType(gin.ErrorTypePrivate).String(),
		}).Info("Handled request")
	}
}
