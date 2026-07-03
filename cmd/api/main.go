package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	configs "github.com/ESE-MONDAY/relay-service/internal/config"
	"github.com/ESE-MONDAY/relay-service/internal/database"
	"github.com/ESE-MONDAY/relay-service/internal/logger"
	"github.com/ESE-MONDAY/relay-service/internal/models"
	"github.com/ESE-MONDAY/relay-service/internal/repository"
	"github.com/ESE-MONDAY/relay-service/internal/router"
)

func main() {

	// Load configuration
	cfg := configs.Load()

	// Initialize logger
	logg, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}
	defer logg.Sync()

	// Connect to PostgreSQL
	dbPool, err := database.NewPool(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	// Initialize repository
	emailRepo := repository.NewEmailRepository(dbPool)

	// Initialize router
	r := router.New(logg)

	// Don't trust all proxies (safe default for local development)
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Temporary endpoint to test the repository
	r.GET("/test-db", func(c *gin.Context) {

		email := &models.Email{
			ID:      uuid.New(),
			From:    "alice@example.com",
			To:      "bob@example.com",
			Subject: "Repository Test",
			Text:    "Hello from Relay Engine",
			HTML:    "<h1>Hello from Relay Engine</h1>",
			Status:  "accepted",
		}

		if err := emailRepo.Save(c.Request.Context(), email); err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "Email saved successfully",
			"id":      email.ID,
		})
	})

	log.Println("🚀 Relay Engine started on :" + cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
