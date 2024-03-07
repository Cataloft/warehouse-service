package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/stdlib"

	"github.com/c9s/goose"

	"lamoda_task/internal/config"
	"lamoda_task/internal/server"
	"lamoda_task/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()
	logger := setupLogger(cfg.Env)

	db := postgres.New(&cfg.Database, logger)

	sqlDB := stdlib.OpenDBFromPool(db.Conn)
	UpMigrations(sqlDB, logger)

	srv := server.New(db, &cfg.Server, logger)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Error("Failed to start server", "error", err)
		}
	}()
	<-done
}

func UpMigrations(db *sql.DB, logger *slog.Logger) {
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("error", "sql migrations", err)
		log.Fatal(err)
	}

	if err := goose.Up(context.Background(), db, "internal/storage/postgres/migrations"); err != nil {
		logger.Error("error", "sql migrations", err)
		log.Fatal(err)
	}
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
