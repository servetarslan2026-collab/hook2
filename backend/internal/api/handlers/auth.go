package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"webhook-service/internal/auth"
	"webhook-service/internal/models"
	"webhook-service/internal/store"
)

type AuthHandler struct {
	store    *store.Store
	authSvc  *auth.AuthService
}

func NewAuthHandler(s *store.Store, a *auth.AuthService) *AuthHandler {
	return &AuthHandler{store: s, authSvc: a}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Email, password, and name are required"})
	}

	if len(req.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Password must be at least 8 characters"})
	}

	// Check if user exists
	existing, _ := h.store.GetUserByEmail(c.Context(), req.Email)
	if existing != nil {
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{Error: "Email already registered"})
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to hash password"})
	}

	user, err := h.store.CreateUser(c.Context(), req.Email, hash, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to create user"})
	}

	token, err := h.authSvc.GenerateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to generate token"})
	}

	refreshToken, err := h.authSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to generate refresh token"})
	}

	return c.Status(fiber.StatusCreated).JSON(models.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Email and password are required"})
	}

	user, err := h.store.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Invalid credentials"})
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Invalid credentials"})
	}

	token, err := h.authSvc.GenerateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to generate token"})
	}

	refreshToken, err := h.authSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to generate refresh token"})
	}

	return c.JSON(models.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req models.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	claims, err := h.authSvc.ValidateToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Invalid refresh token"})
	}

	user, err := h.store.GetUserByID(c.Context(), claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "User not found"})
	}

	token, err := h.authSvc.GenerateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to generate token"})
	}

	refreshToken, err := h.authSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to generate refresh token"})
	}

	return c.JSON(models.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
	})
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req models.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Email is required"})
	}

	// Always return success to prevent email enumeration
	return c.JSON(models.MessageResponse{Message: "If the email exists, a reset link has been sent"})
}

func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
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
