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

	// Email verification
	protected.Post("/auth/send-verification", authH.SendVerificationEmail)
	protected.Post("/auth/verify/:token", authH.VerifyEmail)

	// Organizations
	protected.Post("/organizations", orgH.Create)
	protected.Get("/organizations", orgH.List)
	protected.Get("/organizations/:id", orgH.Get)
	protected.Put("/organizations/:id", orgH.Update)
	protected.Delete("/organizations/:id", orgH.Delete)
	protected.Get("/organizations/:id/dashboard", dashH.OrgStats)
	protected.Get("/organizations/:id/chart", dashH.OrgChartData)

	// Organization members
	protected.Get("/organizations/:id/members", orgH.ListMembers)
	protected.Post("/organizations/:id/members", orgH.InviteMember)
	protected.Delete("/organizations/:id/members/:userId", orgH.RemoveMember)

	// Accept invitation
	protected.Post("/invitations/:token/accept", orgH.AcceptInvitation)

	// Applications
	protected.Post("/organizations/:id/applications", appH.Create)
	protected.Get("/organizations/:id/applications", appH.List)
	protected.Get("/applications/:id", appH.Get)
	protected.Put("/applications/:id", appH.Update)
	protected.Delete("/applications/:id", appH.Delete)
	protected.Get("/applications/:id/dashboard", dashH.AppStats)
	protected.Get("/applications/:id/chart", dashH.AppChartData)

	// Application secrets
	protected.Post("/applications/:id/secrets", secretH.Create)
	protected.Get("/applications/:id/secrets", secretH.List)
	protected.Delete("/applications/:id/secrets/:secretId", secretH.Delete)

	// Event types
	protected.Post("/applications/:id/event-types", eventTypeH.Create)
	protected.Get("/applications/:id/event-types", eventTypeH.List)
	protected.Delete("/applications/:id/event-types/:etId", eventTypeH.Delete)

	// Subscriptions
	protected.Post("/applications/:id/subscriptions", subH.Create)
	protected.Get("/applications/:id/subscriptions", subH.List)
	protected.Get("/subscriptions/:id", subH.Get)
	protected.Put("/subscriptions/:id", subH.Update)
	protected.Delete("/subscriptions/:id", subH.Delete)

	// Events
	protected.Post("/applications/:id/events", eventH.Receive)
	protected.Get("/applications/:id/events", eventH.List)
	protected.Get("/events/:id", eventH.Get)

	// Deliveries
	protected.Get("/applications/:id/deliveries", deliveryH.List)
	protected.Get("/deliveries/:id", deliveryH.Get)
	protected.Post("/deliveries/:id/retry", deliveryH.Retry)

	// API Documentation
	protected.Get("/documentation", handlers.GetAPIDoc())
}
