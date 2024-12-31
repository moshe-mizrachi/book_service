package config

import (
	"book_service/pkg/clients"
	mw "book_service/pkg/middlewares"
	"book_service/pkg/routes"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

func SetupEnv() {
	if os.Getenv("GIN_MODE") == "test" {
		if err := godotenv.Load(".env.test"); err != nil {
			logrus.Infof("Error loading .env.test file: %v", err)
		}
	} else {
		if err := godotenv.Load(); err != nil {
			logrus.Infof("Error loading .env file: %v", err)
		}
	}
}

func SetupServer() *gin.Engine {
	app := gin.New()
	binding.EnableDecoderDisallowUnknownFields = true

	logrus.Infof("Setting up middlewares")
	app.Use(gin.Recovery())
	app.Use(mw.Logger(), mw.RecordActions())

	routes.RegisterRoutes(app)
	logrus.Infof("Middlewares and routes initialized")

	return app
}

func InitClients() {
	if err := clients.InitElasticsearchClient(); err != nil {
		logrus.Infof("Failed to initialize Elasticsearch: %v", err)
	}

	clients.InitRedisClient()
	clients.InitElasticWorkerPool(10)
}
