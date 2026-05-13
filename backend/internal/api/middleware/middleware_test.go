package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestPrometheusMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(PrometheusMiddleware())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestCorsMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(CorsMiddleware("http://localhost:5173"))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
	if resp.Header.Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatalf("Expected CORS header, got: %s", resp.Header.Get("Access-Control-Allow-Origin"))
	}
}
