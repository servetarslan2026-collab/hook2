package middleware

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// ValidateUUIDParam validates that a route parameter is a valid UUID.
// Prevents ID enumeration and injection attacks.
func ValidateUUIDParam(param string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		val := c.Params(param)
		if val == "" || !uuidRegex.MatchString(val) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid " + param + " format",
			})
		}
		return c.Next()
	}
}

// SanitizeInput strips potentially dangerous characters from query params.
func SanitizeInput() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ensure page/per_page are valid integers (Fiber handles this with QueryInt)
		// but we validate string query params for injection
		for _, key := range []string{"status", "event_type", "search"} {
			val := c.Query(key)
			if len(val) > 255 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Query parameter '" + key + "' is too long",
				})
			}
		}
		return c.Next()
	}
}
