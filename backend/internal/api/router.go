package api

import (
	"github.com/gofiber/fiber/v2"
	"webhook-service/internal/api/handlers"
	"webhook-service/internal/api/middleware"
	"webhook-service/internal/auth"
	"webhook-service/internal/queue"
	"webhook-service/internal/store"
)

func SetupRoutes(app *fiber.App, s *store.Store, authSvc *auth.AuthService, q *queue.Queue) {
	// Initialize handlers
	authH := handlers.NewAuthHandler(s, authSvc)
	userH := handlers.NewUserHandler(s, authSvc)
	orgH := handlers.NewOrganizationHandler(s)
	appH := handlers.NewApplicationHandler(s)
	secretH := handlers.NewApplicationSecretHandler(s)
	eventTypeH := handlers.NewEventTypeHandler(s)
	subH := handlers.NewSubscriptionHandler(s)
	eventH := handlers.NewEventHandler(s, q)
	deliveryH := handlers.NewDeliveryHandler(s)
	dashH := handlers.NewDashboardHandler(s)

	api := app.Group("/api/v1")

	// Public routes
	api.Post("/auth/register", authH.Register)
	api.Post("/auth/login", authH.Login)
	api.Post("/auth/refresh", authH.RefreshToken)
	api.Post("/auth/forgot-password", authH.ForgotPassword)

	// Protected routes (JWT or API Key)
	protected := api.Group("", middleware.AuthMiddleware(authSvc, s))

	// User
	protected.Get("/users/me", authH.GetCurrentUser)
	protected.Put("/users/me", userH.UpdateProfile)
	protected.Put("/users/me/password", userH.UpdatePassword)

	// Organizations
	protected.Post("/organizations", orgH.Create)
	protected.Get("/organizations", orgH.List)
	protected.Get("/organizations/:org_id", orgH.Get)
	protected.Put("/organizations/:org_id", orgH.Update)
	protected.Delete("/organizations/:org_id", orgH.Delete)
	protected.Get("/organizations/:org_id/dashboard", dashH.OrgStats)
	protected.Get("/organizations/:org_id/chart", dashH.OrgChartData)

	// Organization members
	protected.Get("/organizations/:org_id/members", orgH.ListMembers)
	protected.Post("/organizations/:org_id/members", orgH.InviteMember)
	protected.Delete("/organizations/:org_id/members/:user_id", orgH.RemoveMember)

	// Applications
	protected.Post("/organizations/:org_id/applications", appH.Create)
	protected.Get("/organizations/:org_id/applications", appH.List)
	protected.Get("/applications/:app_id", appH.Get)
	protected.Put("/applications/:app_id", appH.Update)
	protected.Delete("/applications/:app_id", appH.Delete)
	protected.Get("/applications/:app_id/dashboard", dashH.AppStats)
	protected.Get("/applications/:app_id/chart", dashH.AppChartData)

	// Application secrets
	protected.Post("/applications/:app_id/secrets", secretH.Create)
	protected.Get("/applications/:app_id/secrets", secretH.List)
	protected.Delete("/applications/:app_id/secrets/:secret_id", secretH.Delete)

	// Event types
	protected.Post("/applications/:app_id/event-types", eventTypeH.Create)
	protected.Get("/applications/:app_id/event-types", eventTypeH.List)
	protected.Delete("/applications/:app_id/event-types/:et_id", eventTypeH.Delete)

	// Subscriptions
	protected.Post("/applications/:app_id/subscriptions", subH.Create)
	protected.Get("/applications/:app_id/subscriptions", subH.List)
	protected.Get("/subscriptions/:sub_id", subH.Get)
	protected.Put("/subscriptions/:sub_id", subH.Update)
	protected.Delete("/subscriptions/:sub_id", subH.Delete)

	// Events
	protected.Post("/applications/:app_id/events", eventH.Receive)
	protected.Get("/applications/:app_id/events", eventH.List)
	protected.Get("/events/:event_id", eventH.Get)

	// Deliveries
	protected.Get("/applications/:app_id/deliveries", deliveryH.List)
	protected.Get("/deliveries/:delivery_id", deliveryH.Get)
	protected.Post("/deliveries/:delivery_id/retry", deliveryH.Retry)

	// API Documentation
	protected.Get("/documentation", handlers.GetAPIDoc())
}
