package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/alexzanser/sso/internal/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, port int) *App {

	grpcServer := grpc.NewServer()
	authgrpc.Register(grpcServer)

	return &App{
		log:        log,
		grpcServer: grpc.NewServer(),
		port:       port,
	}
}

func (a *App) MustRun() {
	const op = "grpcapp.MustRun"
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Starting gRPC server", slog.Int("port", a.port))

	if err := a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop gracefully stops the gRPC server.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("Stopping gRPC server", slog.Int("port", a.port))

	a.grpcServer.GracefulStop()
}
