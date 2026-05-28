package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("starting api server")

	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer pool.Close()

	if err := database.RunMigrations(cfg.DatabaseURL, filepath.Join(cfg.BasePath, "migrations")); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	if err := database.LoadSeed(ctx, pool, filepath.Join(cfg.BasePath, "data/seed.json")); err != nil {
		log.Fatal().Err(err).Msg("failed to load seed data")
	}

	router := server.SetupRouter(pool, cfg.JWTSecret, cfg.AllowOrigins)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Str("port", cfg.Port).Msg("server listening")
		if err := router.Run(":" + cfg.Port); err != nil {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	<-quit
	log.Info().Msg("shutting down server")
}
