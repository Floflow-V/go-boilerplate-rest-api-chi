package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "go-boilerplate-rest-api-chi/docs"
	"go-boilerplate-rest-api-chi/internal/api"
	"go-boilerplate-rest-api-chi/internal/config"
	"go-boilerplate-rest-api-chi/internal/database"
	"go-boilerplate-rest-api-chi/internal/logger"
)

// @title						go-boilerplate-rest-api-chi
// @version					1.0
// @description				This is a sample API with Chi.
// @host						localhost:8080
// @BasePath					/api
// @schemes					http
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				JWT security accessToken. Please add it in the format "Bearer {AccessToken}" to authorize your requests.
func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	logger, err := logger.NewLogger(&config)
	if err != nil {
		log.Fatal("failed to init logger", err)
	}

	database, err := database.Init(config, logger)
	if err != nil {
		log.Fatal("failed to init connection with database", err)
	}

	handler := api.CreateApi(config, logger, database.Gorm)

	addr := fmt.Sprintf("%s:%d", config.Api.Host, config.Api.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info().Msgf("Server listening on http://%s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error().Err(err).Msg("Listen error")
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	logger.Info().Msg("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Error().Err(err).Msg("Forced shutdown")
	}

	if err := database.Close(); err != nil {
		logger.Error().Err(err).Msg("Failed to close database")
	}

	logger.Info().Msg("Server and database shutdown cleanly")
}
