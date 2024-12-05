package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "123456789")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)
		log.Print("שדגדש7נגא7שע א77נא7ט ט")

		// access the status we are sendig
		status := c.Writer.Status()
		log.Println(status)
	}
}

func main() {
	r := gin.Default()
	r.Use(Logger())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run("0.0.0.0:1234") // listen and serve on 0.0.0.0:8080
}
