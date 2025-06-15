package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/alexzanser/sso/internal/app/grpc"
	"github.com/alexzanser/sso/internal/services/auth"
	"github.com/alexzanser/sso/internal/storage/sqlite"
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
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.NewService(
		log,
		storage,
		storage,
		storage,
		tokenTTL,
	)

	grpcApp := grpcapp.NewApp(log, grpcPort, authService)
	return &App{
		GRPCsrv: grpcApp,
	}
}
