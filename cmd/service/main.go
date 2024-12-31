package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	app "book_service/pkg/config"
	"book_service/pkg/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	app.SetupEnv()
	app.InitClients()

	server := app.SetupServer()
	port, _ := utils.GetEnvVar[string]("PORT", "8080")
	address := "0.0.0.0:" + port

	httpServer := &http.Server{
		Addr:    address,
		Handler: server.Handler(),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		logrus.Infof("Running server on :%s", port)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Server failed: %s", err)
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logrus.Errorf("Server shutdown failed: %s", err)
	} else {
		logrus.Infof("Server shutdown successfully.")
	}
}
