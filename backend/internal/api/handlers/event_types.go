package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"webhook-service/internal/models"
	"webhook-service/internal/store"
)

type EventTypeHandler struct {
	store *store.Store
}

func NewEventTypeHandler(s *store.Store) *EventTypeHandler {
	return &EventTypeHandler{store: s}
}

func (h *EventTypeHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	appID, err := uuid.Parse(c.Params("appId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid application ID"})
	}

	app, err := h.store.GetApplication(c.Context(), appID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Application not found"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), app.OrganizationID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	var req models.CreateEventTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Name is required"})
	}

	et, err := h.store.CreateEventType(c.Context(), appID, req.Name, req.Description, req.Schema)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to create event type"})
	}

	return c.Status(fiber.StatusCreated).JSON(et)
}

func (h *EventTypeHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	appID, err := uuid.Parse(c.Params("appId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid application ID"})
	}

	app, err := h.store.GetApplication(c.Context(), appID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Application not found"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), app.OrganizationID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	types, total, err := h.store.ListEventTypes(c.Context(), appID, page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to list event types"})
	}

	return c.JSON(models.PaginatedResponse{
		Data:       types,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: total / perPage,
	})
}

func (h *EventTypeHandler) Get(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	etID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid event type ID"})
	}

	et, err := h.store.GetEventType(c.Context(), etID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Event type not found"})
	}

	app, err := h.store.GetApplication(c.Context(), et.ApplicationID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Application not found"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), app.OrganizationID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	return c.JSON(et)
}

func (h *EventTypeHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	etID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid event type ID"})
	}

	et, err := h.store.GetEventType(c.Context(), etID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Event type not found"})
	}

	app, err := h.store.GetApplication(c.Context(), et.ApplicationID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Application not found"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), app.OrganizationID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	var req models.CreateEventTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Name is required"})
	}

	updated, err := h.store.UpdateEventType(c.Context(), etID, req.Name, req.Description, req.Schema)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to update event type"})
	}

	return c.JSON(updated)
}

func (h *EventTypeHandler) Delete(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	etID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid event type ID"})
	}

	et, err := h.store.GetEventType(c.Context(), etID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Event type not found"})
	}

	app, err := h.store.GetApplication(c.Context(), et.ApplicationID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Application not found"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), app.OrganizationID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	if err := h.store.DeleteEventType(c.Context(), etID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to delete event type"})
	}

	return c.JSON(models.MessageResponse{Message: "Event type deleted"})
}
