package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(rdb *redis.Client, limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()

		// Use API key if present, otherwise fall back to IP
		key := "ratelimit:" + c.Get("X-API-Key")
		if key == "ratelimit:" {
			key = "ratelimit:" + c.IP()
		}

		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			return c.Next()
		}

		if count == 1 {
			rdb.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		}

		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", int64(limit)-count))

		return c.Next()
	}
}

// TenantRateLimitMiddleware applies per-organization rate limiting via Redis.
func TenantRateLimitMiddleware(rdb *redis.Client, limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		tenantKey := resolveTenantKey(c)
		if tenantKey == "" {
			return c.Next()
		}

		key := "ratelimit:tenant:" + tenantKey
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			return c.Next()
		}
		if count == 1 {
			rdb.Expire(ctx, key, window)
		}
		if count > int64(limit) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Organization rate limit exceeded. Try again later.",
			})
		}

		c.Set("X-TenantRateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Set("X-TenantRateLimit-Remaining", fmt.Sprintf("%d", int64(limit)-count))
		return c.Next()
	}
}

func resolveTenantKey(c *fiber.Ctx) string {
	// API key auth provides org_id directly
	if orgID := c.Locals("org_id"); orgID != nil {
		if id, ok := orgID.(uuid.UUID); ok {
			return "org:" + id.String()
		}
	}
	// JWT auth — use user_id as proxy
	if userID := c.Locals("user_id"); userID != nil {
		if id, ok := userID.(uuid.UUID); ok {
			return "user:" + id.String()
		}
	}
	return ""
}
