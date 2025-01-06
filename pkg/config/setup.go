package config

import (
	"book_service/pkg/clients"
	"book_service/pkg/consts"
	mw "book_service/pkg/middlewares"
	"book_service/pkg/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func SetupEnv() {
	if os.Getenv("GIN_MODE") == "test" {
		if err := godotenv.Load(".env.test"); err != nil {
			log.Infof("Error loading .env.test file: %v", err)
		}
	} else {
		if err := godotenv.Load(); err != nil {
			log.Infof("Error loading .env file: %v", err)
		}
	}
}

func SetupServer() *gin.Engine {
	app := gin.New()
	binding.EnableDecoderDisallowUnknownFields = true

	log.Infof("Setting up middlewares")
	app.Use(gin.Recovery())
	app.Use(mw.Logger(), mw.RecordActions())

	routes.RegisterRoutes(app)
	log.Infof("Middlewares and routes initialized")

	return app
}

func Setup() *gin.Engine {
	SetupEnv()
	InitClients()
	return SetupServer()
}

func InitClients() {
	if err := clients.InitElasticsearchClient(); err != nil {
		log.Infof("Failed to initialize Elasticsearch: %v", err)
	}

	clients.InitRedisClient()
	clients.InitElasticWorkerPool(consts.WorkersNumber)
}

func ShutDown() {
	ShutDownClients()
}

func ShutDownClients() {
	clients.ShutdownWorkerPool(consts.WorkersNumber)
	clients.ShutDownRedisClient()
}
