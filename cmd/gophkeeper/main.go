package main

import (
	"context"
	"github.com/aseptimu/GophKeeper/internal/app/config"
	"github.com/aseptimu/GophKeeper/internal/app/handlers"
	"github.com/aseptimu/GophKeeper/internal/app/server/http"
	"github.com/aseptimu/GophKeeper/internal/app/services"
	"github.com/aseptimu/GophKeeper/internal/app/store"
	"log/slog"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	appConfig, err := config.NewAppConfig()
	slog.SetLogLoggerLevel(slog.LevelDebug)
	if err != nil {
		slog.Error("Failed to parse config", "error", err)
		return
	}

	slog.Info("Config parsed",
		"Server address", appConfig.ServerAddress)

	if err = config.MigrateDB(appConfig); err != nil {
		slog.Error("Failed to migrate db", "err", err, "db url", appConfig.DSN)
		return
	}

	slog.Info("Starting GophKeeper server", "Address", appConfig.ServerAddress)

	dbStore, err := store.NewDBStore(ctx, appConfig.DSN)
	if err != nil {
		slog.Error("Failed to create db store", "err", err)
		return
	}

	authService := services.NewAuthService(appConfig, dbStore)

	authHandler := handlers.NewAuthHandler(authService)
	http.NewServer(appConfig, authHandler).Run()
}
