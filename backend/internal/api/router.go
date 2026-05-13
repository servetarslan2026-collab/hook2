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
	api := app.Group("/api/v1")

	// Public routes
	api.Post("/auth/register", handlers.Register(s, authSvc))
	api.Post("/auth/login", handlers.Login(s, authSvc))
	api.Post("/auth/refresh", handlers.RefreshToken(authSvc))
	api.Post("/auth/forgot-password", handlers.ForgotPassword(s))

	// Protected routes (JWT or API Key)
	protected := api.Group("", middleware.AuthMiddleware(authSvc, s))

	// User
	protected.Get("/users/me", handlers.GetCurrentUser(s))
	protected.Put("/users/me", handlers.UpdateProfile(s))
	protected.Put("/users/me/password", handlers.UpdatePassword(s, authSvc))

	// Organizations
	protected.Post("/organizations", handlers.CreateOrganization(s))
	protected.Get("/organizations", handlers.ListOrganizations(s))
	protected.Get("/organizations/:org_id", handlers.GetOrganization(s))
	protected.Put("/organizations/:org_id", handlers.UpdateOrganization(s))
	protected.Delete("/organizations/:org_id", handlers.DeleteOrganization(s))
	protected.Get("/organizations/:org_id/dashboard", handlers.OrgDashboard(s))

	// Organization members
	protected.Get("/organizations/:org_id/members", handlers.ListMembers(s))
	protected.Post("/organizations/:org_id/members", handlers.InviteMember(s))
	protected.Delete("/organizations/:org_id/members/:user_id", handlers.RemoveMember(s))

	// Applications
	protected.Post("/organizations/:org_id/applications", handlers.CreateApplication(s))
	protected.Get("/organizations/:org_id/applications", handlers.ListApplications(s))
	protected.Get("/applications/:app_id", handlers.GetApplication(s))
	protected.Put("/applications/:app_id", handlers.UpdateApplication(s))
	protected.Delete("/applications/:app_id", handlers.DeleteApplication(s))
	protected.Get("/applications/:app_id/dashboard", handlers.AppDashboard(s))

	// Application secrets
	protected.Post("/applications/:app_id/secrets", handlers.CreateSecret(s))
	protected.Get("/applications/:app_id/secrets", handlers.ListSecrets(s))
	protected.Delete("/applications/:app_id/secrets/:secret_id", handlers.DeleteSecret(s))

	// Event types
	protected.Post("/applications/:app_id/event-types", handlers.CreateEventType(s))
	protected.Get("/applications/:app_id/event-types", handlers.ListEventTypes(s))
	protected.Delete("/applications/:app_id/event-types/:et_id", handlers.DeleteEventType(s))

	// Subscriptions
	protected.Post("/applications/:app_id/subscriptions", handlers.CreateSubscription(s))
	protected.Get("/applications/:app_id/subscriptions", handlers.ListSubscriptions(s))
	protected.Get("/subscriptions/:sub_id", handlers.GetSubscription(s))
	protected.Put("/subscriptions/:sub_id", handlers.UpdateSubscription(s))
	protected.Delete("/subscriptions/:sub_id", handlers.DeleteSubscription(s))
	protected.Post("/subscriptions/:sub_id/test", handlers.TestSubscription(s, q))

	// Events
	protected.Post("/applications/:app_id/events", handlers.SendEvent(s, q))
	protected.Get("/applications/:app_id/events", handlers.ListEvents(s))
	protected.Get("/events/:event_id", handlers.GetEvent(s))

	// Deliveries
	protected.Get("/applications/:app_id/deliveries", handlers.ListDeliveries(s))
	protected.Get("/deliveries/:delivery_id", handlers.GetDelivery(s))
	protected.Post("/deliveries/:delivery_id/retry", handlers.RetryDelivery(s, q))

	// API Documentation
	protected.Get("/documentation", handlers.GetAPIDoc())
}
