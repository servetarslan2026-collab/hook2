package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	Name          string    `json:"name"`
	IsAdmin       bool      `json:"is_admin"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Organization struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrganizationMember struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	UserID         uuid.UUID `json:"user_id"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
}

type Application struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ApplicationSecret struct {
	ID          uuid.UUID `json:"id"`
	ApplicationID uuid.UUID `json:"application_id"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
}

type EventType struct {
	ID            uuid.UUID       `json:"id"`
	ApplicationID uuid.UUID       `json:"application_id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Schema        json.RawMessage `json:"schema,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type Subscription struct {
	ID            uuid.UUID `json:"id"`
	ApplicationID uuid.UUID `json:"application_id"`
	EventTypes    []string  `json:"event_types"`
	TargetURL     string    `json:"target_url"`
	Secret        string    `json:"secret,omitempty"`
	Description   string    `json:"description"`
	Enabled       bool      `json:"enabled"`
	RetryCount    int       `json:"retry_count"`
	RetryDelays   string    `json:"retry_delays"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Event struct {
	ID            uuid.UUID       `json:"id"`
	ApplicationID uuid.UUID       `json:"application_id"`
	EventType     string          `json:"event_type"`
	Payload       json.RawMessage `json:"payload"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
}

type DeliveryAttempt struct {
	ID             uuid.UUID `json:"id"`
	EventID        uuid.UUID `json:"event_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Status         string    `json:"status"`
	StatusCode     int       `json:"status_code"`
	RequestBody    string    `json:"request_body"`
	ResponseBody   string    `json:"response_body"`
	DurationMs     int64     `json:"duration_ms"`
	AttemptNumber  int       `json:"attempt_number"`
	CreatedAt      time.Time `json:"created_at"`
}

// Request/Response types

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type CreateOrganizationRequest struct {
	Name string `json:"name"`
}

type InviteMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type CreateApplicationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateSecretRequest struct {
	Name string `json:"name"`
}

type CreateEventTypeRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Schema      json.RawMessage `json:"schema,omitempty"`
}

type CreateSubscriptionRequest struct {
	EventTypes  []string `json:"event_types"`
	TargetURL   string   `json:"target_url"`
	Description string   `json:"description"`
	RetryCount  int      `json:"retry_count,omitempty"`
	RetryDelays string   `json:"retry_delays,omitempty"`
}

type SendEventRequest struct {
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
}

type TestDeliveryRequest struct {
	SubscriptionID uuid.UUID `json:"subscription_id"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UpdateOrganizationRequest struct {
	Name string `json:"name"`
}

type UpdateApplicationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateSubscriptionRequest struct {
	EventTypes  *[]string `json:"event_types,omitempty"`
	TargetURL   *string   `json:"target_url,omitempty"`
	Description *string   `json:"description,omitempty"`
	Enabled     *bool     `json:"enabled,omitempty"`
	RetryCount  *int      `json:"retry_count,omitempty"`
	RetryDelays *string   `json:"retry_delays,omitempty"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

type StatsResponse struct {
	TotalEvents      int64   `json:"total_events"`
	TotalDeliveries  int64   `json:"total_deliveries"`
	SuccessRate      float64 `json:"success_rate"`
	TotalSubscriptions int64 `json:"total_subscriptions"`
	TotalApplications  int64 `json:"total_applications"`
}

type ChartDataPoint struct {
	Date    string `json:"date"`
	Success int64  `json:"success"`
	Failed  int64  `json:"failed"`
	Total   int64  `json:"total"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
