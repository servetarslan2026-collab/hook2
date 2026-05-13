package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func GetAPIDoc() fiber.Handler {
	return func(c *fiber.Ctx) error {
		doc := map[string]interface{}{
			"openapi": "3.0.3",
			"info": map[string]string{
				"title":       "Webhook Service API",
				"version":     "1.0.0",
				"description": "Production-ready Webhook-as-a-Service platform",
			},
			"servers": []map[string]string{
				{"url": "/api/v1", "description": "API v1"},
			},
			"paths": map[string]interface{}{
				"/auth/register": map[string]interface{}{
					"post": map[string]interface{}{
						"summary":     "Register a new user",
						"requestBody": map[string]interface{}{"content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"email": map[string]string{"type": "string"}, "password": map[string]string{"type": "string"}, "name": map[string]string{"type": "string"}}}}}},
						"responses":   map[string]interface{}{"201": map[string]string{"description": "User created"}},
					},
				},
				"/auth/login": map[string]interface{}{
					"post": map[string]interface{}{
						"summary":     "Login",
						"requestBody": map[string]interface{}{"content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"email": map[string]string{"type": "string"}, "password": map[string]string{"type": "string"}}}}}},
						"responses":   map[string]interface{}{"200": map[string]string{"description": "Login successful"}},
					},
				},
				"/organizations": map[string]interface{}{
					"get":  map[string]interface{}{"summary": "List organizations", "responses": map[string]interface{}{"200": map[string]string{"description": "Organization list"}}},
					"post": map[string]interface{}{"summary": "Create organization", "responses": map[string]interface{}{"201": map[string]string{"description": "Organization created"}}},
				},
				"/organizations/{org_id}/applications": map[string]interface{}{
					"get":  map[string]interface{}{"summary": "List applications", "responses": map[string]interface{}{"200": map[string]string{"description": "Application list"}}},
					"post": map[string]interface{}{"summary": "Create application", "responses": map[string]interface{}{"201": map[string]string{"description": "Application created"}}},
				},
				"/applications/{app_id}/events": map[string]interface{}{
					"post": map[string]interface{}{
						"summary":     "Send an event (trigger webhooks)",
						"requestBody": map[string]interface{}{"content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"event_type": map[string]string{"type": "string"}, "payload": map[string]string{"type": "object"}}}}}},
						"responses":   map[string]interface{}{"201": map[string]string{"description": "Event created and queued for delivery"}},
					},
					"get": map[string]interface{}{"summary": "List events", "responses": map[string]interface{}{"200": map[string]string{"description": "Event list"}}},
				},
				"/applications/{app_id}/subscriptions": map[string]interface{}{
					"get":  map[string]interface{}{"summary": "List subscriptions", "responses": map[string]interface{}{"200": map[string]string{"description": "Subscription list"}}},
					"post": map[string]interface{}{"summary": "Create subscription", "responses": map[string]interface{}{"201": map[string]string{"description": "Subscription created"}}},
				},
				"/applications/{app_id}/deliveries": map[string]interface{}{
					"get": map[string]interface{}{"summary": "List delivery attempts", "responses": map[string]interface{}{"200": map[string]string{"description": "Delivery list"}}},
				},
			},
		}

		return c.JSON(doc)
	}
}
