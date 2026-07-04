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
	"github.com/ESE-MONDAY/relay-service/internal/queue"
	"github.com/ESE-MONDAY/relay-service/internal/repository"
	"github.com/ESE-MONDAY/relay-service/internal/router"
	"github.com/ESE-MONDAY/relay-service/internal/service"
	"github.com/ESE-MONDAY/relay-service/internal/worker"

	"github.com/ESE-MONDAY/relay-service/internal/processor"
	"github.com/ESE-MONDAY/relay-service/internal/sender"
)

type App struct {
	Config *configs.Config

	Logger *zap.Logger

	DB *pgxpool.Pool

	Server *http.Server

	Workers *worker.Pool

	Context context.Context

	Cancel context.CancelFunc
}

func New() (*App, error) {

	// Configuration
	cfg := configs.Load()

	// Load context
	ctx, cancel := context.WithCancel(
		context.Background(),
	)

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
	jobQueue := queue.NewInMemoryQueue(cfg.QueueSize)

	// Repositories
	emailRepo := repository.NewEmailRepository(dbPool)

	// Services
	emailService := service.NewEmailService(emailRepo, jobQueue)
	emailSender := sender.NewNoopSender()
	emailProcessor := processor.NewEmailProcessor(
		emailRepo,
		emailSender,
		logg,
	)
	//Worker

	workerPool := worker.NewPool(

		cfg.WorkerCount,
		jobQueue,
		emailProcessor,
		logg,
	)
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
		Config:  cfg,
		Logger:  logg,
		DB:      dbPool,
		Server:  server,
		Workers: workerPool,
		Context: ctx,
		Cancel:  cancel,
	}, nil
}

func (a *App) Run() error {

	a.Logger.Info("Starting Relay Engine...")
	a.Workers.Start(a.Context)

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

	// Stop accepting HTTP requests.
	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}

	// Signal all background goroutines to stop.
	a.Cancel()

	// Wait for every worker to finish.
	a.Workers.Wait()

	// Close shared resources.
	a.DB.Close()

	_ = a.Logger.Sync()

	return nil
}
