package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"webhook-service/internal/api"
	"webhook-service/internal/api/middleware"
	"webhook-service/internal/api/ws"
	"webhook-service/internal/auth"
	"webhook-service/internal/config"
	"webhook-service/internal/queue"
	"webhook-service/internal/store"
	"webhook-service/internal/worker"
)

func main() {
	// Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Config
	cfg := config.Load()

	// Database
	ctx := context.Background()
	db, err := pgx.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close(ctx)

	if err := db.Ping(ctx); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}
	logger.Info("Connected to database")

	// Run migrations
	if err := runMigrations(ctx, db); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Store
	s := store.NewStore(db)

	// Auth service
	authSvc := auth.NewAuthService(cfg.JWTSecret, cfg.JWTExpiry)

	// NATS queue
	q, err := queue.Connect(cfg.NATSURL, logger)
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer q.Close()
	logger.Info("Connected to NATS")

	// Redis
	rdb, err := connectRedis(cfg.RedisURL)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer rdb.Close()
	logger.Info("Connected to Redis")

	// WebSocket hub
	hub := ws.NewHub(logger)

	// Start workers
	w := worker.New(q, s, logger, hub)
	if err := w.Start(ctx, cfg.WorkerCount); err != nil {
		logger.Fatal("Failed to start workers", zap.Error(err))
	}

	// Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Webhook Service API",
		BodyLimit:    10 << 20, // 10MB
		ErrorHandler: errorHandler(logger),
	})

	// Middleware
	app.Use(recover.New())
	app.Use(middleware.CorsMiddleware(cfg.FrontendURL))
	app.Use(middleware.LoggerMiddleware(logger))
	app.Use(middleware.PrometheusMiddleware())

	// Prometheus metrics endpoint
	app.Get("/metrics", func(c *fiber.Ctx) error {
		promhttp.Handler().ServeHTTP(c.Context(), nil)
		return nil
	})

	// Routes
	api.SetupRoutes(app, s, authSvc, q, hub)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		logger.Info("Shutting down...")
		app.Shutdown()
	}()

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.Info("Starting server", zap.String("addr", addr))
	if err := app.Listen(addr); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}

func connectRedis(url string) (*redis.Client, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("parse redis URL: %w", err)
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return rdb, nil
}

func runMigrations(ctx context.Context, db *pgx.Conn) error {
	migration := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL DEFAULT '',
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

	CREATE TABLE IF NOT EXISTS organizations (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(255) NOT NULL,
		owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_organizations_owner ON organizations(owner_id);

	CREATE TABLE IF NOT EXISTS organization_members (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		role VARCHAR(50) NOT NULL DEFAULT 'member',
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		UNIQUE(organization_id, user_id)
	);
	CREATE INDEX IF NOT EXISTS idx_org_members_org ON organization_members(organization_id);
	CREATE INDEX IF NOT EXISTS idx_org_members_user ON organization_members(user_id);

	CREATE TABLE IF NOT EXISTS applications (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_applications_org ON applications(organization_id);

	CREATE TABLE IF NOT EXISTS application_secrets (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
		key VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_app_secrets_app ON application_secrets(application_id);
	CREATE INDEX IF NOT EXISTS idx_app_secrets_key ON application_secrets(key);

	CREATE TABLE IF NOT EXISTS event_types (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		schema JSONB,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		UNIQUE(application_id, name)
	);
	CREATE INDEX IF NOT EXISTS idx_event_types_app ON event_types(application_id);

	CREATE TABLE IF NOT EXISTS subscriptions (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
		event_types TEXT[] NOT NULL DEFAULT '{}',
		target_url TEXT NOT NULL,
		secret VARCHAR(255) NOT NULL DEFAULT '',
		description TEXT NOT NULL DEFAULT '',
		enabled BOOLEAN NOT NULL DEFAULT true,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_subscriptions_app ON subscriptions(application_id);
	CREATE INDEX IF NOT EXISTS idx_subscriptions_enabled ON subscriptions(enabled);

	CREATE TABLE IF NOT EXISTS events (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
		event_type VARCHAR(255) NOT NULL,
		payload JSONB NOT NULL DEFAULT '{}',
		metadata JSONB,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_events_app ON events(application_id);
	CREATE INDEX IF NOT EXISTS idx_events_type ON events(event_type);
	CREATE INDEX IF NOT EXISTS idx_events_created ON events(created_at DESC);

	CREATE TABLE IF NOT EXISTS delivery_attempts (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
		subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
		status VARCHAR(50) NOT NULL,
		status_code INT NOT NULL DEFAULT 0,
		request_body TEXT NOT NULL DEFAULT '',
		response_body TEXT NOT NULL DEFAULT '',
		duration_ms BIGINT NOT NULL DEFAULT 0,
		attempt_number INT NOT NULL DEFAULT 1,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_deliveries_event ON delivery_attempts(event_id);
	CREATE INDEX IF NOT EXISTS idx_deliveries_sub ON delivery_attempts(subscription_id);
	CREATE INDEX IF NOT EXISTS idx_deliveries_status ON delivery_attempts(status);
	CREATE INDEX IF NOT EXISTS idx_deliveries_created ON delivery_attempts(created_at DESC);

	CREATE TABLE IF NOT EXISTS invitations (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
		email VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL DEFAULT 'member',
		token VARCHAR(255) UNIQUE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + INTERVAL '7 days'
	);
	CREATE INDEX IF NOT EXISTS idx_invitations_token ON invitations(token);
	CREATE INDEX IF NOT EXISTS idx_invitations_email ON invitations(email);

	CREATE TABLE IF NOT EXISTS verification_tokens (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + INTERVAL '24 hours'
	);
	CREATE INDEX IF NOT EXISTS idx_verification_tokens_token ON verification_tokens(token);

	ALTER TABLE users ADD COLUMN IF NOT EXISTS is_admin BOOLEAN NOT NULL DEFAULT false;
	`

	_, err := db.Exec(ctx, migration)
	return err
}

func errorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		logger.Error("Request error",
			zap.String("path", c.Path()),
			zap.Int("status", code),
			zap.Error(err),
		)
		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
