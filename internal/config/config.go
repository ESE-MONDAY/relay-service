package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppEnv  string

	Port string

	DBHost string
	DBPort string

	DBUser     string
	DBPassword string
	DBName     string

	DatabaseURL string
	WorkerCount int
	QueueSize   int
}

func Load() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}
	workerCount, err := strconv.Atoi(os.Getenv("WORKER_COUNT"))
	if err != nil {
		workerCount = 4 // sensible default
	}

	queueSize, err := strconv.Atoi(os.Getenv("QUEUE_SIZE"))
	if err != nil {
		queueSize = 1000 // sensible default
	}

	return &Config{
		AppName: os.Getenv("APP_NAME"),
		AppEnv:  os.Getenv("APP_ENV"),

		Port: os.Getenv("PORT"),

		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),

		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		WorkerCount: workerCount,
		QueueSize:   queueSize,
	}
}
