package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	ssogrpc "restrocheck/internal/app/grpc/sso"
	"restrocheck/internal/app/rest"
	"restrocheck/internal/config"
	"restrocheck/internal/storage/postgresql"
	"restrocheck/pkg/logger"
	"restrocheck/pkg/logger/sl"
)

func main() {

	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Info("starting logger", slog.String("env", cfg.Env))

	ssoClient, err := ssogrpc.NewClient(
		log,
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)

	if err != nil {
		log.Error("failed to create SSO gRPC client", sl.Err(err))
		os.Exit(1)
	}
	log.Info("SSO gRPC client successfully initialized")
	log.Info("SSO gRPC client created", slog.Any("client", cfg.Clients))

	storage, err := postgresql.NewStorage(cfg)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	log.Info("DataBase success connection")

	app := rest.NewApp(log, storage, ssoClient, cfg)

	go func() {
		app.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.Stop(ctx)
}
