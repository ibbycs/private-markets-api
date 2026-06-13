package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/ibbycs/private-markets-api/internal/config"
	"github.com/ibbycs/private-markets-api/internal/database"
	"github.com/ibbycs/private-markets-api/internal/logger"
	"github.com/ibbycs/private-markets-api/internal/server"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()
	logger := logger.New()

	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseUrl)

	if err != nil {
		logger.Error("An error occurred connecting to the database", "error", err)
		return
	}

	defer pool.Close()

	err = database.Migrate(logger, pool)

	if err != nil {
		logger.Error("An error occurred during migrate", "error", err)
		return
	}

	e, err := server.NewServer(cfg, logger, pool)

	go func() {
		logger.Info("Open API Documentation: http://localhost:" + cfg.HostPort + "/docs")
		err := e.Start(":" + cfg.HostPort)

		if err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info("Server gracefully stopped")
}
