package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"webhook-service/internal/models"
	"webhook-service/internal/store"
)

type SubscriptionHandler struct {
	store *store.Store
}

func NewSubscriptionHandler(s *store.Store) *SubscriptionHandler {
	return &SubscriptionHandler{store: s}
}

func (h *SubscriptionHandler) Create(c *fiber.Ctx) error {
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

	var req models.CreateSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if req.TargetURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Target URL is required"})
	}

	if len(req.EventTypes) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "At least one event type is required"})
	}

	sub, err := h.store.CreateSubscription(c.Context(), appID, req.EventTypes, req.TargetURL, req.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to create subscription"})
	}

	return c.Status(fiber.StatusCreated).JSON(sub)
}

func (h *SubscriptionHandler) List(c *fiber.Ctx) error {
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

	subs, total, err := h.store.ListSubscriptions(c.Context(), appID, page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to list subscriptions"})
	}

	return c.JSON(models.PaginatedResponse{
		Data:       subs,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: total / perPage,
	})
}

func (h *SubscriptionHandler) Get(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	subID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid subscription ID"})
	}

	sub, err := h.store.GetSubscription(c.Context(), subID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Subscription not found"})
	}

	app, err := h.store.GetApplication(c.Context(), sub.ApplicationID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Application not found"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), app.OrganizationID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	return c.JSON(sub)
}

func (h *SubscriptionHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	subID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid subscription ID"})
	}

	sub, err := h.store.GetSubscription(c.Context(), subID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Subscription not found"})
	}

	app, err := h.store.GetApplication(c.Context(), sub.ApplicationID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Application not found"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), app.OrganizationID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	var req models.UpdateSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	updated, err := h.store.UpdateSubscription(c.Context(), subID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to update subscription"})
	}

	return c.JSON(updated)
}

func (h *SubscriptionHandler) Delete(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	subID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid subscription ID"})
	}

	sub, err := h.store.GetSubscription(c.Context(), subID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Subscription not found"})
	}

	app, err := h.store.GetApplication(c.Context(), sub.ApplicationID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Application not found"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), app.OrganizationID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	if err := h.store.DeleteSubscription(c.Context(), subID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to delete subscription"})
	}

	return c.JSON(models.MessageResponse{Message: "Subscription deleted"})
}
