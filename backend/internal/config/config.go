package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds process settings loaded from the environment.
// Keep names boring and explicit so docker-compose and .env.example stay readable.
type Config struct {
	HTTPAddr       string
	MetricsAddr    string
	DatabaseURL    string
	RedisURL       string
	AWSRegion      string
	AWSAccessKey   string
	AWSSecretKey   string
	SQSEndpoint    string
	SQSQueueURL    string
	PublicAppURL   string
	SeedPassword   string
	SessionCookie  string
	SessionTTL     time.Duration
	AWSEndpointURL string // optional LocalStack override
}

func Load() (Config, error) {
	ttlHours, err := strconv.Atoi(getenv("SESSION_TTL_HOURS", "168"))
	if err != nil {
		return Config{}, fmt.Errorf("SESSION_TTL_HOURS: %w", err)
	}

	cfg := Config{
		HTTPAddr:       getenv("HTTP_ADDR", ":8080"),
		MetricsAddr:    getenv("METRICS_ADDR", ":8081"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		RedisURL:       getenv("REDIS_URL", "redis://localhost:6379/0"),
		AWSRegion:      getenv("AWS_REGION", "us-east-1"),
		AWSAccessKey:   getenv("AWS_ACCESS_KEY_ID", "x"),
		AWSSecretKey:   getenv("AWS_SECRET_ACCESS_KEY", "x"),
		SQSEndpoint:    os.Getenv("SQS_ENDPOINT"),
		SQSQueueURL:    os.Getenv("SQS_QUEUE_URL"),
		PublicAppURL:   getenv("PUBLIC_APP_URL", "http://localhost:5173"),
		SeedPassword:   getenv("SEED_PASSWORD", "phishguard-dev"),
		SessionCookie:  getenv("SESSION_COOKIE_NAME", "phishguard_session"),
		SessionTTL:     time.Duration(ttlHours) * time.Hour,
		AWSEndpointURL: os.Getenv("AWS_ENDPOINT_URL"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.SQSQueueURL == "" {
		return Config{}, fmt.Errorf("SQS_QUEUE_URL is required")
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
