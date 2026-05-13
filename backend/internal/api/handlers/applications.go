package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"webhook-service/internal/models"
	"webhook-service/internal/store"
)

type ApplicationHandler struct {
	store *store.Store
}

func NewApplicationHandler(s *store.Store) *ApplicationHandler {
	return &ApplicationHandler{store: s}
}

func (h *ApplicationHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	orgID, err := uuid.Parse(c.Params("orgId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid organization ID"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), orgID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	var req models.CreateApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Name is required"})
	}

	app, err := h.store.CreateApplication(c.Context(), orgID, req.Name, req.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to create application"})
	}

	return c.Status(fiber.StatusCreated).JSON(app)
}

func (h *ApplicationHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	orgID, err := uuid.Parse(c.Params("orgId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid organization ID"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), orgID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	apps, total, err := h.store.ListApplications(c.Context(), orgID, page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to list applications"})
	}

	return c.JSON(models.PaginatedResponse{
		Data:       apps,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: total / perPage,
	})
}

func (h *ApplicationHandler) Get(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	appID, err := uuid.Parse(c.Params("id"))
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

	return c.JSON(app)
}

func (h *ApplicationHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	appID, err := uuid.Parse(c.Params("id"))
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

	var req models.UpdateApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Name is required"})
	}

	updated, err := h.store.UpdateApplication(c.Context(), appID, req.Name, req.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to update application"})
	}

	return c.JSON(updated)
}

func (h *ApplicationHandler) Delete(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	appID, err := uuid.Parse(c.Params("id"))
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

	if err := h.store.DeleteApplication(c.Context(), appID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to delete application"})
	}

	return c.JSON(models.MessageResponse{Message: "Application deleted"})
}
