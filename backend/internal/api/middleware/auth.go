package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"webhook-service/internal/auth"
	"webhook-service/internal/store"
)

func AuthMiddleware(authSvc *auth.AuthService, store *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try JWT first
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := authSvc.ValidateToken(token)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
			}
			c.Locals("user_id", claims.UserID)
			c.Locals("email", claims.Email)
			c.Locals("auth_type", "jwt")
			return c.Next()
		}

		// Try API key
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}
		if apiKey != "" {
			app, err := store.GetApplicationByAPIKey(c.Context(), apiKey)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid API key"})
			}
			c.Locals("app_id", app.ID)
			c.Locals("org_id", app.OrganizationID)
			c.Locals("auth_type", "api_key")
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication required"})
	}
}
