package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func GetAPIDoc() fiber.Handler {
	return func(c *fiber.Ctx) error {
		doc := map[string]interface{}{
			"openapi": "3.0.3",
			"info": map[string]interface{}{
				"title":       "Webhook Service API",
				"version":     "1.0.0",
				"description": "Production-ready Webhook-as-a-Service platform. Send events, deliver webhooks, manage subscriptions.",
			},
			"servers": []map[string]string{
				{"url": "/api/v1", "description": "API v1"},
			},
			"components": map[string]interface{}{
				"securitySchemes": map[string]interface{}{
					"BearerAuth": map[string]string{"type": "http", "scheme": "bearer", "bearerFormat": "JWT"},
					"ApiKeyAuth":  map[string]string{"type": "apiKey", "in": "header", "name": "X-API-Key"},
				},
				"schemas": map[string]interface{}{
					"Error": map[string]interface{}{
						"type": "object", "properties": map[string]interface{}{
							"error": map[string]string{"type": "string"},
						},
					},
					"User": map[string]interface{}{
						"type": "object", "properties": map[string]interface{}{
							"id": map[string]string{"type": "string", "format": "uuid"},
							"email": map[string]string{"type": "string"},
							"name": map[string]string{"type": "string"},
							"email_verified": map[string]string{"type": "boolean"},
							"created_at": map[string]string{"type": "string", "format": "date-time"},
						},
					},
					"Organization": map[string]interface{}{
						"type": "object", "properties": map[string]interface{}{
							"id": map[string]string{"type": "string", "format": "uuid"},
							"name": map[string]string{"type": "string"},
							"owner_id": map[string]string{"type": "string", "format": "uuid"},
							"created_at": map[string]string{"type": "string", "format": "date-time"},
						},
					},
					"Application": map[string]interface{}{
						"type": "object", "properties": map[string]interface{}{
							"id": map[string]string{"type": "string", "format": "uuid"},
							"organization_id": map[string]string{"type": "string", "format": "uuid"},
							"name": map[string]string{"type": "string"},
							"description": map[string]string{"type": "string"},
							"created_at": map[string]string{"type": "string", "format": "date-time"},
						},
					},
					"Subscription": map[string]interface{}{
						"type": "object", "properties": map[string]interface{}{
							"id": map[string]string{"type": "string", "format": "uuid"},
							"application_id": map[string]string{"type": "string", "format": "uuid"},
							"event_types": map[string]interface{}{"type": "array", "items": map[string]string{"type": "string"}},
							"target_url": map[string]string{"type": "string", "format": "uri"},
							"description": map[string]string{"type": "string"},
							"enabled": map[string]string{"type": "boolean"},
							"created_at": map[string]string{"type": "string", "format": "date-time"},
						},
					},
					"Event": map[string]interface{}{
						"type": "object", "properties": map[string]interface{}{
							"id": map[string]string{"type": "string", "format": "uuid"},
							"application_id": map[string]string{"type": "string", "format": "uuid"},
							"event_type": map[string]string{"type": "string"},
							"payload": map[string]string{"type": "object"},
							"metadata": map[string]string{"type": "object"},
							"created_at": map[string]string{"type": "string", "format": "date-time"},
						},
					},
					"DeliveryAttempt": map[string]interface{}{
						"type": "object", "properties": map[string]interface{}{
							"id": map[string]string{"type": "string", "format": "uuid"},
							"event_id": map[string]string{"type": "string", "format": "uuid"},
							"subscription_id": map[string]string{"type": "string", "format": "uuid"},
							"status": map[string]string{"type": "string", "enum": "success,failed,dead_letter,pending"},
							"status_code": map[string]string{"type": "integer"},
							"request_body": map[string]string{"type": "string"},
							"response_body": map[string]string{"type": "string"},
							"duration_ms": map[string]string{"type": "integer"},
							"attempt_number": map[string]string{"type": "integer"},
							"created_at": map[string]string{"type": "string", "format": "date-time"},
						},
					},
				},
			},
			"security": []map[string]interface{}{
				{"BearerAuth": []string{}},
				{"ApiKeyAuth": []string{}},
			},
			"paths": map[string]interface{}{
				// ===== Auth =====
				"/auth/register": map[string]interface{}{
					"post": map[string]interface{}{
						"tags":        []string{"Auth"},
						"summary":     "Register a new user",
						"security":    []map[string]interface{}{},
						"requestBody": map[string]interface{}{"content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "required": []string{"email", "password", "name"}, "properties": map[string]interface{}{"email": map[string]string{"type": "string", "format": "email"}, "password": map[string]string{"type": "string", "minLength": "8"}, "name": map[string]string{"type": "string"}}}}}},
						"responses":   map[string]interface{}{"201": map[string]interface{}{"description": "User created", "content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"token": map[string]string{"type": "string"}, "refresh_token": map[string]string{"type": "string"}, "user": map[string]interface{}{"$ref": "#/components/schemas/User"}}}}}}},
					},
				},
				"/auth/login": map[string]interface{}{
					"post": map[string]interface{}{
						"tags":        []string{"Auth"},
						"summary":     "Login",
						"security":    []map[string]interface{}{},
						"requestBody": map[string]interface{}{"content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "required": []string{"email", "password"}, "properties": map[string]interface{}{"email": map[string]string{"type": "string"}, "password": map[string]string{"type": "string"}}}}}},
						"responses":   map[string]interface{}{"200": map[string]interface{}{"description": "Login successful"}},
					},
				},
				"/auth/send-verification": map[string]interface{}{
					"post": map[string]interface{}{"tags": []string{"Auth"}, "summary": "Send verification email", "responses": map[string]interface{}{"200": map[string]string{"description": "Verification email sent"}}},
				},
				"/auth/verify/{token}": map[string]interface{}{
					"post": map[string]interface{}{"tags": []string{"Auth"}, "summary": "Verify email with token", "parameters": []map[string]interface{}{{"name": "token", "in": "path", "required": true, "schema": map[string]string{"type": "string"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Email verified"}}},
				},
				// ===== Users =====
				"/users/me": map[string]interface{}{
					"get":  map[string]interface{}{"tags": []string{"Users"}, "summary": "Get current user", "responses": map[string]interface{}{"200": map[string]interface{}{"description": "User profile", "content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"$ref": "#/components/schemas/User"}}}}}},
					"put": map[string]interface{}{"tags": []string{"Users"}, "summary": "Update profile", "requestBody": map[string]interface{}{"content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"name": map[string]string{"type": "string"}}}}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Profile updated"}}},
				},
				"/users/me/password": map[string]interface{}{
					"put": map[string]interface{}{"tags": []string{"Users"}, "summary": "Change password", "requestBody": map[string]interface{}{"content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"current_password": map[string]string{"type": "string"}, "new_password": map[string]string{"type": "string", "minLength": "8"}}}}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Password changed"}}},
				},
				// ===== Organizations =====
				"/organizations": map[string]interface{}{
					"get":  map[string]interface{}{"tags": []string{"Organizations"}, "summary": "List organizations", "responses": map[string]interface{}{"200": map[string]string{"description": "Organization list"}}},
					"post": map[string]interface{}{"tags": []string{"Organizations"}, "summary": "Create organization", "requestBody": map[string]interface{}{"content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "required": []string{"name"}, "properties": map[string]interface{}{"name": map[string]string{"type": "string"}}}}}}, "responses": map[string]interface{}{"201": map[string]string{"description": "Organization created"}}},
				},
				"/organizations/{id}": map[string]interface{}{
					"get":    map[string]interface{}{"tags": []string{"Organizations"}, "summary": "Get organization", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Organization details"}}},
					"put":    map[string]interface{}{"tags": []string{"Organizations"}, "summary": "Update organization", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Organization updated"}}},
					"delete": map[string]interface{}{"tags": []string{"Organizations"}, "summary": "Delete organization", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Organization deleted"}}},
				},
				"/organizations/{id}/members": map[string]interface{}{
					"get":  map[string]interface{}{"tags": []string{"Organizations"}, "summary": "List members", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Member list"}}},
					"post": map[string]interface{}{"tags": []string{"Organizations"}, "summary": "Invite member", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "requestBody": map[string]interface{}{"content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"email": map[string]string{"type": "string"}, "role": map[string]string{"type": "string"}}}}}}, "responses": map[string]interface{}{"201": map[string]string{"description": "Member added or invitation created"}}},
				},
				"/organizations/{id}/members/{userId}": map[string]interface{}{
					"delete": map[string]interface{}{"tags": []string{"Organizations"}, "summary": "Remove member", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}, {"name": "userId", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Member removed"}}},
				},
				"/invitations/{token}/accept": map[string]interface{}{
					"post": map[string]interface{}{"tags": []string{"Organizations"}, "summary": "Accept invitation", "parameters": []map[string]interface{}{{"name": "token", "in": "path", "required": true, "schema": map[string]string{"type": "string"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Invitation accepted"}}},
				},
				"/organizations/{id}/dashboard": map[string]interface{}{
					"get": map[string]interface{}{"tags": []string{"Dashboard"}, "summary": "Organization stats", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Stats response"}}},
				},
				// ===== Applications =====
				"/organizations/{orgId}/applications": map[string]interface{}{
					"get":  map[string]interface{}{"tags": []string{"Applications"}, "summary": "List applications", "parameters": []map[string]interface{}{{"name": "orgId", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Application list"}}},
					"post": map[string]interface{}{"tags": []string{"Applications"}, "summary": "Create application", "parameters": []map[string]interface{}{{"name": "orgId", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "requestBody": map[string]interface{}{"content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "required": []string{"name"}, "properties": map[string]interface{}{"name": map[string]string{"type": "string"}, "description": map[string]string{"type": "string"}}}}}}, "responses": map[string]interface{}{"201": map[string]string{"description": "Application created"}}},
				},
				"/applications/{id}": map[string]interface{}{
					"get":    map[string]interface{}{"tags": []string{"Applications"}, "summary": "Get application", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Application details"}}},
					"put":    map[string]interface{}{"tags": []string{"Applications"}, "summary": "Update application", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Application updated"}}},
					"delete": map[string]interface{}{"tags": []string{"Applications"}, "summary": "Delete application", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Application deleted"}}},
				},
				"/applications/{id}/dashboard": map[string]interface{}{
					"get": map[string]interface{}{"tags": []string{"Dashboard"}, "summary": "Application stats + chart", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Stats + chart data"}}},
				},
				// ===== Secrets =====
				"/applications/{id}/secrets": map[string]interface{}{
					"get":  map[string]interface{}{"tags": []string{"Secrets"}, "summary": "List API keys", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}, {"name": "page", "in": "query", "schema": map[string]string{"type": "integer"}}, {"name": "per_page", "in": "query", "schema": map[string]string{"type": "integer"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Paginated secret list"}}},
					"post": map[string]interface{}{"tags": []string{"Secrets"}, "summary": "Create API key", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "requestBody": map[string]interface{}{"content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "required": []string{"name"}, "properties": map[string]interface{}{"name": map[string]string{"type": "string"}}}}}}, "responses": map[string]interface{}{"201": map[string]string{"description": "Secret created"}}},
				},
				"/applications/{id}/secrets/{secretId}": map[string]interface{}{
					"delete": map[string]interface{}{"tags": []string{"Secrets"}, "summary": "Delete API key", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}, {"name": "secretId", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Secret deleted"}}},
				},
				// ===== Event Types =====
				"/applications/{id}/event-types": map[string]interface{}{
					"get":  map[string]interface{}{"tags": []string{"Event Types"}, "summary": "List event types", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}, {"name": "page", "in": "query", "schema": map[string]string{"type": "integer"}}, {"name": "per_page", "in": "query", "schema": map[string]string{"type": "integer"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Paginated event type list"}}},
					"post": map[string]interface{}{"tags": []string{"Event Types"}, "summary": "Create event type", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "requestBody": map[string]interface{}{"content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "required": []string{"name"}, "properties": map[string]interface{}{"name": map[string]string{"type": "string"}, "description": map[string]string{"type": "string"}, "schema": map[string]string{"type": "object"}}}}}}, "responses": map[string]interface{}{"201": map[string]string{"description": "Event type created"}}},
				},
				"/applications/{id}/event-types/{etId}": map[string]interface{}{
					"delete": map[string]interface{}{"tags": []string{"Event Types"}, "summary": "Delete event type", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}, {"name": "etId", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Event type deleted"}}},
				},
				// ===== Subscriptions =====
				"/applications/{id}/subscriptions": map[string]interface{}{
					"get":  map[string]interface{}{"tags": []string{"Subscriptions"}, "summary": "List subscriptions", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}, {"name": "page", "in": "query", "schema": map[string]string{"type": "integer"}}, {"name": "per_page", "in": "query", "schema": map[string]string{"type": "integer"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Paginated subscription list"}}},
					"post": map[string]interface{}{"tags": []string{"Subscriptions"}, "summary": "Create subscription", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "requestBody": map[string]interface{}{"content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "required": []string{"event_types", "target_url"}, "properties": map[string]interface{}{"event_types": map[string]interface{}{"type": "array", "items": map[string]string{"type": "string"}}, "target_url": map[string]string{"type": "string", "format": "uri"}, "description": map[string]string{"type": "string"}}}}}}, "responses": map[string]interface{}{"201": map[string]string{"description": "Subscription created"}}},
				},
				"/subscriptions/{id}": map[string]interface{}{
					"get":    map[string]interface{}{"tags": []string{"Subscriptions"}, "summary": "Get subscription", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Subscription details"}}},
					"put":    map[string]interface{}{"tags": []string{"Subscriptions"}, "summary": "Update subscription", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Subscription updated"}}},
					"delete": map[string]interface{}{"tags": []string{"Subscriptions"}, "summary": "Delete subscription", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Subscription deleted"}}},
				},
				// ===== Events =====
				"/applications/{id}/events": map[string]interface{}{
					"post": map[string]interface{}{
						"tags":     []string{"Events"},
						"summary":  "Send event (trigger webhooks)",
						"parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}},
						"requestBody": map[string]interface{}{"content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "required": []string{"event_type", "payload"}, "properties": map[string]interface{}{"event_type": map[string]string{"type": "string"}, "payload": map[string]string{"type": "object"}, "metadata": map[string]string{"type": "object"}}}}}},
						"responses": map[string]interface{}{"201": map[string]interface{}{"description": "Event created", "content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"$ref": "#/components/schemas/Event"}}}}},
					},
					"get": map[string]interface{}{"tags": []string{"Events"}, "summary": "List events", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}, {"name": "page", "in": "query", "schema": map[string]string{"type": "integer"}}, {"name": "per_page", "in": "query", "schema": map[string]string{"type": "integer"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Paginated event list"}}},
				},
				"/events/{id}": map[string]interface{}{
					"get": map[string]interface{}{"tags": []string{"Events"}, "summary": "Get event details", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]interface{}{"description": "Event details", "content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"$ref": "#/components/schemas/Event"}}}}}},
				},
				// ===== Deliveries =====
				"/applications/{id}/deliveries": map[string]interface{}{
					"get": map[string]interface{}{"tags": []string{"Deliveries"}, "summary": "List delivery attempts", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}, {"name": "page", "in": "query", "schema": map[string]string{"type": "integer"}}, {"name": "per_page", "in": "query", "schema": map[string]string{"type": "integer"}}, {"name": "status", "in": "query", "schema": map[string]string{"type": "string"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Paginated delivery list"}}},
				},
				"/deliveries/{id}": map[string]interface{}{
					"get": map[string]interface{}{"tags": []string{"Deliveries"}, "summary": "Get delivery details", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]interface{}{"description": "Delivery details", "content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"$ref": "#/components/schemas/DeliveryAttempt"}}}}}},
				},
				"/deliveries/{id}/retry": map[string]interface{}{
					"post": map[string]interface{}{"tags": []string{"Deliveries"}, "summary": "Retry delivery", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Delivery retried"}}},
				},
				// ===== Health =====
				"/health": map[string]interface{}{
					"get": map[string]interface{}{"tags": []string{"System"}, "summary": "Health check", "security": []map[string]interface{}{}, "responses": map[string]interface{}{"200": map[string]interface{}{"description": "Service healthy", "content": map[string]interface{}{"application/json": map[string]interface{}{"schema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"status": map[string]string{"type": "string"}}}}}}}},
				},
			},
		}

		c.Set("Content-Type", "application/json")
		return c.JSON(doc)
	}
}
