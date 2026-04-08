package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/pg")
	t.Setenv("SQS_QUEUE_URL", "http://localhost/queue")
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("METRICS_ADDR", "")
	t.Setenv("REDIS_URL", "")
	t.Setenv("AWS_REGION", "")
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")
	t.Setenv("PUBLIC_APP_URL", "")
	t.Setenv("SEED_PASSWORD", "")
	t.Setenv("SESSION_COOKIE_NAME", "")
	t.Setenv("SESSION_TTL_HOURS", "")
	t.Setenv("SQS_ENDPOINT", "")
	t.Setenv("AWS_ENDPOINT_URL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("HTTPAddr=%q", cfg.HTTPAddr)
	}
	if cfg.MetricsAddr != ":8081" {
		t.Fatalf("MetricsAddr=%q", cfg.MetricsAddr)
	}
	if cfg.RedisURL != "redis://localhost:6379/0" {
		t.Fatalf("RedisURL=%q", cfg.RedisURL)
	}
	if cfg.AWSRegion != "us-east-1" {
		t.Fatalf("AWSRegion=%q", cfg.AWSRegion)
	}
	if cfg.PublicAppURL != "http://localhost:5173" {
		t.Fatalf("PublicAppURL=%q", cfg.PublicAppURL)
	}
	if cfg.SessionCookie != "phishguard_session" {
		t.Fatalf("SessionCookie=%q", cfg.SessionCookie)
	}
	if cfg.SessionTTL != 168*time.Hour {
		t.Fatalf("SessionTTL=%v", cfg.SessionTTL)
	}
	if cfg.SeedPassword != "phishguard-dev" {
		t.Fatalf("SeedPassword=%q", cfg.SeedPassword)
	}
}

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("SQS_QUEUE_URL", "http://localhost/queue")
	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Fatalf("err=%v", err)
	}
}

func TestLoadRequiresSQSQueueURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/pg")
	t.Setenv("SQS_QUEUE_URL", "")
	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "SQS_QUEUE_URL") {
		t.Fatalf("err=%v", err)
	}
}

func TestLoadBadSessionTTL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/pg")
	t.Setenv("SQS_QUEUE_URL", "http://localhost/queue")
	t.Setenv("SESSION_TTL_HOURS", "not-a-number")
	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "SESSION_TTL_HOURS") {
		t.Fatalf("err=%v", err)
	}
}
