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

func Setup() *gin.Engine {
	setupEnv()
	initClients()
	return setupServer()
}

func ShutDown() {
	shutDownClients()
}

func setupEnv() {
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

func setupServer() *gin.Engine {
	app := gin.New()
	binding.EnableDecoderDisallowUnknownFields = true

	log.Infof("Setting up middlewares")
	app.Use(gin.Recovery())
	app.Use(mw.Logger(), mw.RecordActions())

	routes.RegisterRoutes(app)
	log.Infof("Middlewares and routes initialized")

	return app
}

func initClients() {
	if err := clients.InitElasticsearchClient(); err != nil {
		log.Infof("Failed to initialize Elasticsearch: %v", err)
	}

	clients.InitRedisClient()
	clients.InitElasticWorkerPool(consts.WorkersNumber)
}

func shutDownClients() {
	clients.ShutdownWorkerPool(consts.WorkersNumber)
	clients.ShutDownRedisClient()
}
