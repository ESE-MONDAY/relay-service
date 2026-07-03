package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	configs "github.com/ESE-MONDAY/relay-service/internal/config"
	"github.com/ESE-MONDAY/relay-service/internal/database"
	"github.com/ESE-MONDAY/relay-service/internal/handler"
	"github.com/ESE-MONDAY/relay-service/internal/logger"
	"github.com/ESE-MONDAY/relay-service/internal/repository"
	"github.com/ESE-MONDAY/relay-service/internal/router"
	"github.com/ESE-MONDAY/relay-service/internal/service"
)

type App struct {
	Config *configs.Config

	Logger *zap.Logger

	DB *pgxpool.Pool

	Server *http.Server
}

func New() (*App, error) {

	// Configuration
	cfg := configs.Load()

	// Logger
	logg, err := logger.New()
	if err != nil {
		return nil, err
	}

	// Database
	dbPool, err := database.NewPool(cfg)
	if err != nil {
		return nil, err
	}

	// Repositories
	emailRepo := repository.NewEmailRepository(dbPool)

	// Services
	emailService := service.NewEmailService(emailRepo)

	// Handlers
	handlers := handler.New(emailService)

	// Router
	r := router.New(logg)

	router.Register(r, handlers)

	// HTTP Server
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		Config: cfg,
		Logger: logg,
		DB:     dbPool,
		Server: server,
	}, nil
}

func (a *App) Run() error {

	a.Logger.Info("Starting Relay Engine...")

	err := a.Server.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (a *App) Shutdown(
	ctx context.Context,
) error {

	a.Logger.Info("Shutting down Relay Engine...")

	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}

	// Close the database pool.
	a.DB.Close()

	// Flush any buffered log entries.
	_ = a.Logger.Sync()

	return nil
}
