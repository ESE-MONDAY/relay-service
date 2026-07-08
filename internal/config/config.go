package configs

import (
	"log"
	"os"
	"strconv"
	"strings"

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

	// Legacy (can be removed after old worker system is deleted)
	WorkerCount int
	QueueSize   int

	// SMTP
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	// Kafka / Redpanda
	KafkaBrokers    []string
	KafkaTopic      string
	KafkaGroupID    string
	KafkaRetryTopic string
	KafkaDLQTopic   string

	MaxRetries int
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	workerCount, err := strconv.Atoi(os.Getenv("WORKER_COUNT"))
	if err != nil {
		workerCount = 4
	}

	queueSize, err := strconv.Atoi(os.Getenv("QUEUE_SIZE"))
	if err != nil {
		queueSize = 1000
	}

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		smtpPort = 2525
	}

	maxRetries, err := strconv.Atoi(os.Getenv("MAX_RETRIES"))
	if err != nil {
		maxRetries = 3
	}

	// Parse broker list.
	// Example:
	// KAFKA_BROKERS=localhost:19092
	// KAFKA_BROKERS=broker1:9092,broker2:9092
	var brokers []string

	if value := strings.TrimSpace(os.Getenv("KAFKA_BROKERS")); value != "" {
		for _, broker := range strings.Split(value, ",") {
			broker = strings.TrimSpace(broker)
			if broker != "" {
				brokers = append(brokers, broker)
			}
		}
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

		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     smtpPort,
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:     os.Getenv("SMTP_FROM"),

		KafkaBrokers:    brokers,
		KafkaTopic:      os.Getenv("KAFKA_TOPIC"),
		KafkaGroupID:    os.Getenv("KAFKA_GROUP_ID"),
		KafkaRetryTopic: os.Getenv("KAFKA_RETRY_TOPIC"),
		KafkaDLQTopic:   os.Getenv("KAFKA_DLQ_TOPIC"),

		MaxRetries: maxRetries,
	}
}
