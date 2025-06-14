package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/alexzanser/sso/internal/app/grpc"
)

type App struct {
	GRPCsrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	grpcApp := grpcapp.NewApp(log, grpcPort)
	return &App{
		GRPCsrv: grpcApp,
	}
}
