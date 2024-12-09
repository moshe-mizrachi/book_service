package main

import (
	"book_service/pkg/clients"
	"book_service/pkg/routes"
	"book_service/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}

	if err := clients.InitElasticsearchClient(); err != nil {
		log.Fatalf("Failed to initialize Elasticsearch: %v", err)
	}
	clients.InitWorkerPool(10)

	r := gin.New()
	binding.EnableDecoderDisallowUnknownFields = true

	r.Use(gin.Recovery())
	r.Use(utils.CustomLogger())

	routes.RegisterRoutes(r)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run("0.0.0.0:1234") // listen and serve on 0.0.0.0:8080
}
