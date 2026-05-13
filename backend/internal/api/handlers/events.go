package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"webhook-service/internal/models"
	"webhook-service/internal/queue"
	"webhook-service/internal/store"
)

type EventHandler struct {
	store *store.Store
	queue *queue.Queue
}

func NewEventHandler(s *store.Store, q *queue.Queue) *EventHandler {
	return &EventHandler{store: s, queue: q}
}

func (h *EventHandler) Receive(c *fiber.Ctx) error {
	appIDVal := c.Locals("app_id")
	if appIDVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Authentication required"})
	}
	appID := appIDVal.(uuid.UUID)

	var req models.SendEventRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if req.EventType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Event type is required"})
	}

	if req.Payload == nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Payload is required"})
	}

	// Validate JSON
	if !json.Valid(req.Payload) {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid JSON payload"})
	}

	if req.Metadata == nil {
		req.Metadata = json.RawMessage("{}")
	}

	event, err := h.store.CreateEvent(c.Context(), appID, req.EventType, req.Payload, req.Metadata)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to create event"})
	}

	// Publish to queue for delivery
	if err := h.queue.PublishEvent(event); err != nil {
		fmt.Printf("Warning: failed to publish event to queue: %v\n", err)
	}

	return c.Status(fiber.StatusCreated).JSON(event)
}

func (h *EventHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		// API key auth
		appIDVal := c.Locals("app_id")
		if appIDVal == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Authentication required"})
		}
		appID := appIDVal.(uuid.UUID)
		eventType := c.Query("event_type", "")
		page := c.QueryInt("page", 1)
		perPage := c.QueryInt("per_page", 20)

		events, total, err := h.store.ListEvents(c.Context(), appID, eventType, page, perPage)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to list events"})
		}

		return c.JSON(models.PaginatedResponse{
			Data:       events,
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: total / perPage,
		})
	}

	appID, err := uuid.Parse(c.Params("appId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid application ID"})
	}

	app, err := h.store.GetApplication(c.Context(), appID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Application not found"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), app.OrganizationID, userID.(uuid.UUID))
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	eventType := c.Query("event_type", "")
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	events, total, err := h.store.ListEvents(c.Context(), appID, eventType, page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to list events"})
	}

	return c.JSON(models.PaginatedResponse{
		Data:       events,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: total / perPage,
	})
}

func (h *EventHandler) Get(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid event ID"})
	}

	event, err := h.store.GetEvent(c.Context(), eventID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Event not found"})
	}

	if userID != nil {
		app, err := h.store.GetApplication(c.Context(), event.ApplicationID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Application not found"})
		}
		isMember, err := h.store.IsOrganizationMember(c.Context(), app.OrganizationID, userID.(uuid.UUID))
		if err != nil || !isMember {
			return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
		}
	}

	return c.JSON(event)
}
