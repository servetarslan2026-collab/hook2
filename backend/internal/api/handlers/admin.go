package handlers

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"webhook-service/internal/models"
	"webhook-service/internal/store"
)

func totalPages(count int, perPage int) int {
	return int(math.Ceil(float64(count) / float64(perPage)))
}

type AdminHandler struct {
	store *store.Store
}

func NewAdminHandler(s *store.Store) *AdminHandler {
	return &AdminHandler{store: s}
}

// SystemStats returns overall system statistics.
func (h *AdminHandler) SystemStats(c *fiber.Ctx) error {
	stats, err := h.store.GetSystemStats(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to get system stats"})
	}
	return c.JSON(stats)
}

// ListUsers returns all users with pagination.
func (h *AdminHandler) ListUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	users, total, err := h.store.ListAllUsers(c.Context(), page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to list users"})
	}

	return c.JSON(models.PaginatedResponse{
		Data:       users,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages(total, perPage),
	})
}

// SetUserAdmin toggles admin status for a user.
func (h *AdminHandler) SetUserAdmin(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid user ID"})
	}

	var req struct {
		IsAdmin bool `json:"is_admin"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	// Prevent self-demotion
	currentUserID := c.Locals("user_id").(uuid.UUID)
	if userID == currentUserID && !req.IsAdmin {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Cannot remove your own admin status"})
	}

	if err := h.store.SetUserAdmin(c.Context(), userID, req.IsAdmin); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to update user"})
	}

	return c.JSON(models.MessageResponse{Message: "User admin status updated"})
}

// DeleteUser deletes a user (admin only).
func (h *AdminHandler) DeleteUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid user ID"})
	}

	// Prevent self-deletion
	currentUserID := c.Locals("user_id").(uuid.UUID)
	if userID == currentUserID {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Cannot delete your own account"})
	}

	if err := h.store.DeleteUser(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to delete user"})
	}

	return c.JSON(models.MessageResponse{Message: "User deleted"})
}

// ListOrganizations returns all organizations (admin view).
func (h *AdminHandler) ListOrganizations(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	orgs, total, err := h.store.ListAllOrganizations(c.Context(), page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to list organizations"})
	}

	return c.JSON(models.PaginatedResponse{
		Data:       orgs,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages(total, perPage),
	})
}

// DeleteOrganization deletes an organization (admin only).
func (h *AdminHandler) DeleteOrganization(c *fiber.Ctx) error {
	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid organization ID"})
	}

	if err := h.store.DeleteOrganization(c.Context(), orgID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to delete organization"})
	}

	return c.JSON(models.MessageResponse{Message: "Organization deleted"})
}

// ListDeadLetters returns dead letter queue entries.
func (h *AdminHandler) ListDeadLetters(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	attempts, total, err := h.store.ListDeadLetterDeliveries(c.Context(), page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to list dead letters"})
	}

	return c.JSON(models.PaginatedResponse{
		Data:       attempts,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages(total, perPage),
	})
}
