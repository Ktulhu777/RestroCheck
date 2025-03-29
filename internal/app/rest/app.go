package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	ssogrpc "restrocheck/internal/app/grpc/sso"
	"restrocheck/internal/config"
	"restrocheck/internal/repository"
	"restrocheck/internal/service"
	"restrocheck/internal/storage"
	"restrocheck/internal/transport/rest/handlers"
	mwJWTAuth "restrocheck/internal/transport/rest/middleware/authentication/jwt"

	mwLogger "restrocheck/internal/transport/rest/middleware/logger"
	"restrocheck/pkg/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type App struct {
	log     *slog.Logger
	server  *http.Server
	storage storage.Storage
	cfg     *config.Config
}

func NewApp(log *slog.Logger, storage *storage.Storage, ssoapi *ssogrpc.Client, cfg *config.Config) *App {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Recoverer)
	router.Use(mwLogger.New(log))
	
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}))

	repos := repository.NewRepositories(storage.DB)
	service := service.NewService(service.Deps{
		Repos: repos,
	})
	
	waiterHandler := handlers.NewWaiterHandler(log, service.Waiter)
	categoryHandler := handlers.NewCategoryHandler(log, service.Category)
	menuHandler := handlers.NewMenuHandler(log, service.Menu)
	priceHandler := handlers.NewPriceHandler(log, service.Price)
	orderHandler := handlers.NewOrderHandler(log, service.Order)
	authJWTMiddleware := mwJWTAuth.JWTAuthIsAdminMiddleware(log, ssoapi, cfg.AppSecret)

	router.Route("/", func(r chi.Router) {
		r.Use(authJWTMiddleware)

		// waiter handlers
		r.Post("/waiter", waiterHandler.SaveWaiter())
		r.Get("/waiter/{id}", waiterHandler.FetchWaiter())
		r.Patch("/waiter/{id}", waiterHandler.ChangeWaiter())
		r.Delete("/waiter/{id}", waiterHandler.RemoveWaiter())
		r.Get("/waiters", waiterHandler.FetchAllWaiters())

		// category handlers
		r.Post("/category", categoryHandler.SaveCategory())

		// menu handlers
		r.Post("/menu", menuHandler.SaveMenu())

		// price handlers
		r.Post("/price", priceHandler.SavePrice())

		// order handlers
		r.Post("/order", orderHandler.SaveOrder())
	})
	

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	return &App{
		log:     log,
		server:  srv,
		storage: *storage,
		cfg:     cfg,
	}
}

func (a *App) Run() error {
	const fn = "internal.app.rest.Run"

	a.log.Info("starting http server", slog.String("address", a.server.Addr))

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Stop(ctx context.Context) {
	const fn = "internal.app.rest.Stop"

	a.log.Info("stopping http server", slog.String("address", a.server.Addr))

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.log.Error("failed to stop server", slog.String("fn", fn), sl.Err(err))
	} else {
		a.log.Info("http server stopped")
	}

	if err := a.storage.Close(); err != nil {
		a.log.Error("failed to close database", slog.Any("error", err))
	} else {
		a.log.Info("database connection closed")
	}
}
