package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/ESE-MONDAY/relay-service/internal/app"
)

func main() {

	// -------------------------------------------------------------------------
	// Bootstrap application
	// -------------------------------------------------------------------------

	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	// -------------------------------------------------------------------------
	// Root context
	// -------------------------------------------------------------------------

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// -------------------------------------------------------------------------
	// errgroup manages all long-running goroutines
	// -------------------------------------------------------------------------

	g, ctx := errgroup.WithContext(ctx)

	// -------------------------------------------------------------------------
	// HTTP Server
	// -------------------------------------------------------------------------

	g.Go(func() error {
		return application.Run()
	})

	// -------------------------------------------------------------------------
	// Wait for shutdown signal
	// -------------------------------------------------------------------------

	<-ctx.Done()

	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	// -------------------------------------------------------------------------
	// Wait for all goroutines
	// -------------------------------------------------------------------------

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	log.Println("Relay Engine stopped")
}
