package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"book_service/pkg/config"
	"book_service/pkg/utils"

	log "github.com/sirupsen/logrus"
)

func main() {
	server := config.Setup()

	port, _ := utils.GetEnvVar[string]("PORT", "8080")
	address := "0.0.0.0:" + port

	httpServer := &http.Server{
		Addr:    address,
		Handler: server.Handler(),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		log.Infof("Running server on :%s", port)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %s", err)
		}
	}()

	<-quit
	config.ShutDown()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Errorf("Server shutdown failed: %s", err)
	} else {
		log.Infof("Server shutdown successfully.")
	}
}
