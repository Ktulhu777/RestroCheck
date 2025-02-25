package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	storage, err := postgresql.NewStorage(cfg)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	log.Info("DataBase success connection")

	app := rest.NewApp(log, storage, cfg.Address, cfg.Timeout, cfg.IdleTimeout)

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