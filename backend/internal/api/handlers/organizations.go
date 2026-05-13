package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"webhook-service/internal/models"
	"webhook-service/internal/store"
)

type OrganizationHandler struct {
	store *store.Store
}

func NewOrganizationHandler(s *store.Store) *OrganizationHandler {
	return &OrganizationHandler{store: s}
}

func (h *OrganizationHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	var req models.CreateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Name is required"})
	}

	org, err := h.store.CreateOrganization(c.Context(), req.Name, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to create organization"})
	}

	return c.Status(fiber.StatusCreated).JSON(org)
}

func (h *OrganizationHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	orgs, total, err := h.store.ListOrganizationsByUser(c.Context(), userID, page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to list organizations"})
	}

	return c.JSON(models.PaginatedResponse{
		Data:       orgs,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: total / perPage,
	})
}

func (h *OrganizationHandler) Get(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid organization ID"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), orgID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	org, err := h.store.GetOrganization(c.Context(), orgID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Organization not found"})
	}

	return c.JSON(org)
}

func (h *OrganizationHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid organization ID"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), orgID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	var req models.UpdateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Name is required"})
	}

	org, err := h.store.UpdateOrganization(c.Context(), orgID, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to update organization"})
	}

	return c.JSON(org)
}

func (h *OrganizationHandler) Delete(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid organization ID"})
	}

	org, err := h.store.GetOrganization(c.Context(), orgID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "Organization not found"})
	}

	if org.OwnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Only the owner can delete an organization"})
	}

	if err := h.store.DeleteOrganization(c.Context(), orgID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to delete organization"})
	}

	return c.JSON(models.MessageResponse{Message: "Organization deleted"})
}

func (h *OrganizationHandler) ListMembers(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid organization ID"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), orgID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	members, total, err := h.store.ListMembersWithUserDetails(c.Context(), orgID, page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to list members"})
	}

	return c.JSON(models.PaginatedResponse{
		Data:       members,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: total / perPage,
	})
}

func (h *OrganizationHandler) InviteMember(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid organization ID"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), orgID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	var req models.InviteMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Email is required"})
	}

	if req.Role == "" {
		req.Role = "member"
	}

	// Try to find user by email
	invitedUser, err := h.store.GetUserByEmail(c.Context(), req.Email)
	if err == nil {
		// User exists, add them directly
		member, err := h.store.AddOrganizationMember(c.Context(), orgID, invitedUser.ID, req.Role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to add member"})
		}
		return c.Status(fiber.StatusCreated).JSON(member)
	}

	// User doesn't exist, create invitation
	inv, err := h.store.CreateInvitation(c.Context(), orgID, req.Email, req.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to create invitation"})
	}

	return c.Status(fiber.StatusCreated).JSON(inv)
}

func (h *OrganizationHandler) RemoveMember(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Not authenticated"})
	}

	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid organization ID"})
	}

	memberID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid user ID"})
	}

	isMember, err := h.store.IsOrganizationMember(c.Context(), orgID, userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Error: "Access denied"})
	}

	if err := h.store.RemoveOrganizationMember(c.Context(), orgID, memberID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to remove member"})
	}

	return c.JSON(models.MessageResponse{Message: "Member removed"})
}
