package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	NATSURL     string
	RedisURL    string
	JWTSecret   string
	JWTExpiry   time.Duration
	WorkerCount int
	FrontendURL string
}

func Load() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "3000"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://webhook:webhook@localhost:5432/webhook_service?sslmode=disable"),
		NATSURL:     getEnv("NATS_URL", "nats://localhost:4222"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production-please"),
		JWTExpiry:   getDuration("JWT_EXPIRY", 24*time.Hour),
		WorkerCount: getInt("WORKER_COUNT", 4),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
