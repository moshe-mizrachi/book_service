package main

import (
	app "book_service/pkg"
	"book_service/pkg/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	app.SetupEnv()
	app.InitClients()

	server := app.SetupServer()
	port, _ := utils.GetEnvVar[string]("PORT", "8080")
	logrus.Infof("Running server on :%s", port)
	server.Run("0.0.0.0:" + port)
}
