package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	configs "github.com/ESE-MONDAY/relay-service/internal/config"
	"github.com/ESE-MONDAY/relay-service/internal/consumer"
	"github.com/ESE-MONDAY/relay-service/internal/database"
	"github.com/ESE-MONDAY/relay-service/internal/handler"
	"github.com/ESE-MONDAY/relay-service/internal/logger"
	"github.com/ESE-MONDAY/relay-service/internal/processor"
	"github.com/ESE-MONDAY/relay-service/internal/queue"
	"github.com/ESE-MONDAY/relay-service/internal/repository"
	"github.com/ESE-MONDAY/relay-service/internal/router"
	"github.com/ESE-MONDAY/relay-service/internal/sender"
	"github.com/ESE-MONDAY/relay-service/internal/service"
)

type App struct {
	Config *configs.Config
	Logger *zap.Logger
	DB     *pgxpool.Pool

	Server *http.Server

	Consumer *consumer.EmailConsumer

	Context context.Context
	Cancel  context.CancelFunc
}

func New() (*App, error) {

	// Configuration
	cfg := configs.Load()

	ctx, cancel := context.WithCancel(context.Background())

	// Logger
	logg, err := logger.New()
	if err != nil {
		cancel()
		return nil, err
	}

	// Database
	dbPool, err := database.NewPool(cfg)
	if err != nil {
		cancel()
		return nil, err
	}

	// Event publisher (Kafka / Redpanda)
	publisher := queue.NewKafkaPublisher(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
	)

	// Repository
	emailRepo := repository.NewEmailRepository(dbPool)

	// SMTP sender
	smtpSender := sender.NewSMTPSender(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.SMTPFrom,
	)

	// Email processor
	emailProcessor := processor.NewEmailProcessor(
		emailRepo,
		publisher,
		smtpSender,
		logg,
	)

	// Consumer
	emailConsumer := consumer.NewEmailConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		cfg.KafkaGroupID,
		emailProcessor,
	)

	go emailConsumer.Start(ctx)

	// Service
	emailService := service.NewEmailService(
		emailRepo,
		publisher,
	)

	// HTTP handlers
	handlers := handler.New(emailService)

	// Router
	r := router.New(logg)
	router.Register(r, handlers)

	// HTTP server
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		Config:   cfg,
		Logger:   logg,
		DB:       dbPool,
		Server:   server,
		Consumer: emailConsumer,
		Context:  ctx,
		Cancel:   cancel,
	}, nil
}

func (a *App) Run() error {

	a.Logger.Info("Starting Relay Engine...")

	if err := a.Server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {

	a.Logger.Info("Shutting down Relay Engine...")

	// Stop accepting HTTP requests.
	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}

	// Cancel background goroutines.
	a.Cancel()

	// Close database pool.
	a.DB.Close()

	// Flush logger.
	_ = a.Logger.Sync()

	return nil
}
