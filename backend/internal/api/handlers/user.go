package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"webhook-service/internal/auth"
	"webhook-service/internal/models"
	"webhook-service/internal/store"
)

type UserHandler struct {
	store   *store.Store
	authSvc *auth.AuthService
}

func NewUserHandler(s *store.Store, a *auth.AuthService) *UserHandler {
	return &UserHandler{store: s, authSvc: a}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	user, err := h.store.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "User not found"})
	}

	return c.JSON(user)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	var req models.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if req.Name == "" || req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Name and email are required"})
	}

	user, err := h.store.UpdateUser(c.Context(), userID, req.Name, req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to update profile"})
	}

	return c.JSON(user)
}

func (h *UserHandler) UpdatePassword(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	var req models.UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Current and new password are required"})
	}

	if len(req.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "New password must be at least 8 characters"})
	}

	user, err := h.store.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "User not found"})
	}

	if !auth.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Current password is incorrect"})
	}

	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to hash password"})
	}

	if err := h.store.UpdateUserPassword(c.Context(), userID, hash); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to update password"})
	}

	return c.JSON(models.MessageResponse{Message: "Password updated"})
}

func (h *UserHandler) ListOrganizations(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	orgs, err := h.store.GetUserOrganizations(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to list organizations"})
	}

	return c.JSON(orgs)
}
