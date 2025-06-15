package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexzanser/sso/internal/app"
	"github.com/alexzanser/sso/internal/config"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("app started", slog.Any("config", cfg))

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	application.GRPCsrv.MustRun()

	
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	application.GRPCsrv.Stop()
	log.Info("app stopped")
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
