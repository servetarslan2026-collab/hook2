package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"webhook-service/internal/models"
	"webhook-service/internal/store"
)

type DashboardHandler struct {
	store *store.Store
}

func NewDashboardHandler(s *store.Store) *DashboardHandler {
	return &DashboardHandler{store: s}
}

func (h *DashboardHandler) AppStats(c *fiber.Ctx) error {
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

	stats, err := h.store.GetAppStats(c.Context(), appID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to get stats"})
	}

	return c.JSON(stats)
}

func (h *DashboardHandler) AppChartData(c *fiber.Ctx) error {
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

	days := c.QueryInt("days", 30)

	data, err := h.store.GetChartData(c.Context(), appID, days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to get chart data"})
	}

	return c.JSON(data)
}

func (h *DashboardHandler) OrgStats(c *fiber.Ctx) error {
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

	stats, err := h.store.GetOrgStats(c.Context(), orgID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to get stats"})
	}

	return c.JSON(stats)
}

func (h *DashboardHandler) OrgChartData(c *fiber.Ctx) error {
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

	days := c.QueryInt("days", 30)

	// Get all apps in org
	apps, _, err := h.store.ListApplications(c.Context(), orgID, 1, 100)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to get chart data"})
	}

	// Aggregate chart data across all apps
	allData := make(map[string]*models.ChartDataPoint)
	for _, app := range apps {
		data, err := h.store.GetChartData(c.Context(), app.ID, days)
		if err != nil {
			continue
		}
		for _, d := range data {
			if existing, ok := allData[d.Date]; ok {
				existing.Success += d.Success
				existing.Failed += d.Failed
				existing.Total += d.Total
			} else {
				allData[d.Date] = &models.ChartDataPoint{
					Date:    d.Date,
					Success: d.Success,
					Failed:  d.Failed,
					Total:   d.Total,
				}
			}
		}
	}

	var result []models.ChartDataPoint
	for _, v := range allData {
		result = append(result, *v)
	}

	return c.JSON(result)
}
