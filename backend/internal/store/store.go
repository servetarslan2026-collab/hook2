package store

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"webhook-service/internal/models"
)

type Store struct {
	db *pgx.Conn
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{db: db}
}

// Pagination helper
func paginate(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return page, perPage
}

func totalPages(count int, perPage int) int {
	return int(math.Ceil(float64(count) / float64(perPage)))
}

// ==================== Users ====================

func (s *Store) CreateUser(ctx context.Context, email, passwordHash, name string) (*models.User, error) {
	u := &models.User{}
	err := s.db.QueryRow(ctx,
		`INSERT INTO users (id, email, password_hash, name, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $5)
		 RETURNING id, email, name, created_at, updated_at`,
		uuid.New(), email, passwordHash, name, time.Now(),
	).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

func (s *Store) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	u := &models.User{}
	err := s.db.QueryRow(ctx,
		`SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return u, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	err := s.db.QueryRow(ctx,
		`SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return u, nil
}

func (s *Store) UpdateUser(ctx context.Context, id uuid.UUID, name, email string) (*models.User, error) {
	u := &models.User{}
	err := s.db.QueryRow(ctx,
		`UPDATE users SET name = $2, email = $3, updated_at = $4 WHERE id = $1
		 RETURNING id, email, name, created_at, updated_at`,
		id, name, email, time.Now(),
	).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return u, nil
}

func (s *Store) UpdateUserPassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	_, err := s.db.Exec(ctx,
		`UPDATE users SET password_hash = $2, updated_at = $3 WHERE id = $1`,
		id, passwordHash, time.Now(),
	)
	return err
}

// ==================== Organizations ====================

func (s *Store) CreateOrganization(ctx context.Context, name string, ownerID uuid.UUID) (*models.Organization, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	o := &models.Organization{}
	err = tx.QueryRow(ctx,
		`INSERT INTO organizations (id, name, owner_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $4)
		 RETURNING id, name, owner_id, created_at, updated_at`,
		uuid.New(), name, ownerID, time.Now(),
	).Scan(&o.ID, &o.Name, &o.OwnerID, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO organization_members (id, organization_id, user_id, role, created_at)
		 VALUES ($1, $2, $3, 'owner', $4)`,
		uuid.New(), o.ID, ownerID, time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("add owner as member: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Store) GetOrganization(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	o := &models.Organization{}
	err := s.db.QueryRow(ctx,
		`SELECT id, name, owner_id, created_at, updated_at FROM organizations WHERE id = $1`, id,
	).Scan(&o.ID, &o.Name, &o.OwnerID, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}
	return o, nil
}

func (s *Store) ListOrganizationsByUser(ctx context.Context, userID uuid.UUID, page, perPage int) ([]models.Organization, int, error) {
	page, perPage = paginate(page, perPage)
	offset := (page - 1) * perPage

	var total int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM organizations o
		 JOIN organization_members om ON o.id = om.organization_id
		 WHERE om.user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(ctx,
		`SELECT o.id, o.name, o.owner_id, o.created_at, o.updated_at
		 FROM organizations o
		 JOIN organization_members om ON o.id = om.organization_id
		 WHERE om.user_id = $1
		 ORDER BY o.created_at DESC
		 LIMIT $2 OFFSET $3`, userID, perPage, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orgs []models.Organization
	for rows.Next() {
		var o models.Organization
		if err := rows.Scan(&o.ID, &o.Name, &o.OwnerID, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, err
		}
		orgs = append(orgs, o)
	}
	return orgs, total, nil
}

func (s *Store) UpdateOrganization(ctx context.Context, id uuid.UUID, name string) (*models.Organization, error) {
	o := &models.Organization{}
	err := s.db.QueryRow(ctx,
		`UPDATE organizations SET name = $2, updated_at = $3 WHERE id = $1
		 RETURNING id, name, owner_id, created_at, updated_at`,
		id, name, time.Now(),
	).Scan(&o.ID, &o.Name, &o.OwnerID, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("update organization: %w", err)
	}
	return o, nil
}

func (s *Store) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `DELETE FROM organizations WHERE id = $1`, id)
	return err
}

func (s *Store) IsOrganizationMember(ctx context.Context, orgID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM organization_members WHERE organization_id = $1 AND user_id = $2)`,
		orgID, userID,
	).Scan(&exists)
	return exists, err
}

// ==================== Organization Members ====================

func (s *Store) AddOrganizationMember(ctx context.Context, orgID, userID uuid.UUID, role string) (*models.OrganizationMember, error) {
	m := &models.OrganizationMember{}
	err := s.db.QueryRow(ctx,
		`INSERT INTO organization_members (id, organization_id, user_id, role, created_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, organization_id, user_id, role, created_at`,
		uuid.New(), orgID, userID, role, time.Now(),
	).Scan(&m.ID, &m.OrganizationID, &m.UserID, &m.Role, &m.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("add member: %w", err)
	}
	return m, nil
}

func (s *Store) ListOrganizationMembers(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]models.OrganizationMember, int, error) {
	page, perPage = paginate(page, perPage)
	offset := (page - 1) * perPage

	var total int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM organization_members WHERE organization_id = $1`, orgID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(ctx,
		`SELECT id, organization_id, user_id, role, created_at
		 FROM organization_members WHERE organization_id = $1
		 ORDER BY created_at ASC LIMIT $2 OFFSET $3`, orgID, perPage, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var members []models.OrganizationMember
	for rows.Next() {
		var m models.OrganizationMember
		if err := rows.Scan(&m.ID, &m.OrganizationID, &m.UserID, &m.Role, &m.CreatedAt); err != nil {
			return nil, 0, err
		}
		members = append(members, m)
	}
	return members, total, nil
}

func (s *Store) RemoveOrganizationMember(ctx context.Context, orgID, userID uuid.UUID) error {
	_, err := s.db.Exec(ctx,
		`DELETE FROM organization_members WHERE organization_id = $1 AND user_id = $2`,
		orgID, userID,
	)
	return err
}

func (s *Store) UpdateMemberRole(ctx context.Context, orgID, userID uuid.UUID, role string) error {
	_, err := s.db.Exec(ctx,
		`UPDATE organization_members SET role = $3 WHERE organization_id = $1 AND user_id = $2`,
		orgID, userID, role,
	)
	return err
}

// ==================== Applications ====================

func (s *Store) CreateApplication(ctx context.Context, orgID uuid.UUID, name, description string) (*models.Application, error) {
	a := &models.Application{}
	err := s.db.QueryRow(ctx,
		`INSERT INTO applications (id, organization_id, name, description, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $5)
		 RETURNING id, organization_id, name, description, created_at, updated_at`,
		uuid.New(), orgID, name, description, time.Now(),
	).Scan(&a.ID, &a.OrganizationID, &a.Name, &a.Description, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create application: %w", err)
	}
	return a, nil
}

func (s *Store) GetApplication(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	a := &models.Application{}
	err := s.db.QueryRow(ctx,
		`SELECT id, organization_id, name, description, created_at, updated_at FROM applications WHERE id = $1`, id,
	).Scan(&a.ID, &a.OrganizationID, &a.Name, &a.Description, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get application: %w", err)
	}
	return a, nil
}

func (s *Store) ListApplications(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]models.Application, int, error) {
	page, perPage = paginate(page, perPage)
	offset := (page - 1) * perPage

	var total int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM applications WHERE organization_id = $1`, orgID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(ctx,
		`SELECT id, organization_id, name, description, created_at, updated_at
		 FROM applications WHERE organization_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, orgID, perPage, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var apps []models.Application
	for rows.Next() {
		var a models.Application
		if err := rows.Scan(&a.ID, &a.OrganizationID, &a.Name, &a.Description, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, 0, err
		}
		apps = append(apps, a)
	}
	return apps, total, nil
}

func (s *Store) UpdateApplication(ctx context.Context, id uuid.UUID, name, description string) (*models.Application, error) {
	a := &models.Application{}
	err := s.db.QueryRow(ctx,
		`UPDATE applications SET name = $2, description = $3, updated_at = $4 WHERE id = $1
		 RETURNING id, organization_id, name, description, created_at, updated_at`,
		id, name, description, time.Now(),
	).Scan(&a.ID, &a.OrganizationID, &a.Name, &a.Description, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("update application: %w", err)
	}
	return a, nil
}

func (s *Store) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `DELETE FROM applications WHERE id = $1`, id)
	return err
}

// ==================== Application Secrets ====================

func (s *Store) CreateApplicationSecret(ctx context.Context, appID uuid.UUID, name string) (*models.ApplicationSecret, error) {
	sec := &models.ApplicationSecret{}
	key := "whsec_" + uuid.New().String()
	err := s.db.QueryRow(ctx,
		`INSERT INTO application_secrets (id, application_id, key, name, created_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, application_id, key, name, created_at`,
		uuid.New(), appID, key, name, time.Now(),
	).Scan(&sec.ID, &sec.ApplicationID, &sec.Key, &sec.Name, &sec.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create secret: %w", err)
	}
	return sec, nil
}

func (s *Store) GetApplicationSecretByKey(ctx context.Context, key string) (*models.ApplicationSecret, error) {
	sec := &models.ApplicationSecret{}
	err := s.db.QueryRow(ctx,
		`SELECT id, application_id, key, name, created_at FROM application_secrets WHERE key = $1`, key,
	).Scan(&sec.ID, &sec.ApplicationID, &sec.Key, &sec.Name, &sec.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get secret by key: %w", err)
	}
	return sec, nil
}

func (s *Store) ListApplicationSecrets(ctx context.Context, appID uuid.UUID, page, perPage int) ([]models.ApplicationSecret, int, error) {
	page, perPage = paginate(page, perPage)
	offset := (page - 1) * perPage

	var total int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM application_secrets WHERE application_id = $1`, appID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(ctx,
		`SELECT id, application_id, key, name, created_at
		 FROM application_secrets WHERE application_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, appID, perPage, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var secrets []models.ApplicationSecret
	for rows.Next() {
		var sec models.ApplicationSecret
		if err := rows.Scan(&sec.ID, &sec.ApplicationID, &sec.Key, &sec.Name, &sec.CreatedAt); err != nil {
			return nil, 0, err
		}
		secrets = append(secrets, sec)
	}
	return secrets, total, nil
}

func (s *Store) DeleteApplicationSecret(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `DELETE FROM application_secrets WHERE id = $1`, id)
	return err
}

// ==================== Event Types ====================

func (s *Store) CreateEventType(ctx context.Context, appID uuid.UUID, name, description string, schema json.RawMessage) (*models.EventType, error) {
	et := &models.EventType{}
	var schemaVal pgtype.Text
	if schema != nil {
		schemaVal = pgtype.Text{String: string(schema), Valid: true}
	}
	err := s.db.QueryRow(ctx,
		`INSERT INTO event_types (id, application_id, name, description, schema, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $6)
		 RETURNING id, application_id, name, description, schema, created_at, updated_at`,
		uuid.New(), appID, name, description, schemaVal, time.Now(),
	).Scan(&et.ID, &et.ApplicationID, &et.Name, &et.Description, &schemaVal, &et.CreatedAt, &et.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create event type: %w", err)
	}
	if schemaVal.Valid {
		et.Schema = json.RawMessage(schemaVal.String)
	}
	return et, nil
}

func (s *Store) GetEventType(ctx context.Context, id uuid.UUID) (*models.EventType, error) {
	et := &models.EventType{}
	var schemaVal pgtype.Text
	err := s.db.QueryRow(ctx,
		`SELECT id, application_id, name, description, schema, created_at, updated_at FROM event_types WHERE id = $1`, id,
	).Scan(&et.ID, &et.ApplicationID, &et.Name, &et.Description, &schemaVal, &et.CreatedAt, &et.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get event type: %w", err)
	}
	if schemaVal.Valid {
		et.Schema = json.RawMessage(schemaVal.String)
	}
	return et, nil
}

func (s *Store) ListEventTypes(ctx context.Context, appID uuid.UUID, page, perPage int) ([]models.EventType, int, error) {
	page, perPage = paginate(page, perPage)
	offset := (page - 1) * perPage

	var total int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM event_types WHERE application_id = $1`, appID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(ctx,
		`SELECT id, application_id, name, description, schema, created_at, updated_at
		 FROM event_types WHERE application_id = $1
		 ORDER BY name ASC LIMIT $2 OFFSET $3`, appID, perPage, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var types []models.EventType
	for rows.Next() {
		var et models.EventType
		var schemaVal pgtype.Text
		if err := rows.Scan(&et.ID, &et.ApplicationID, &et.Name, &et.Description, &schemaVal, &et.CreatedAt, &et.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if schemaVal.Valid {
			et.Schema = json.RawMessage(schemaVal.String)
		}
		types = append(types, et)
	}
	return types, total, nil
}

func (s *Store) UpdateEventType(ctx context.Context, id uuid.UUID, name, description string, schema json.RawMessage) (*models.EventType, error) {
	et := &models.EventType{}
	var schemaVal pgtype.Text
	if schema != nil {
		schemaVal = pgtype.Text{String: string(schema), Valid: true}
	}
	err := s.db.QueryRow(ctx,
		`UPDATE event_types SET name = $2, description = $3, schema = $4, updated_at = $5 WHERE id = $1
		 RETURNING id, application_id, name, description, schema, created_at, updated_at`,
		id, name, description, schemaVal, time.Now(),
	).Scan(&et.ID, &et.ApplicationID, &et.Name, &et.Description, &schemaVal, &et.CreatedAt, &et.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("update event type: %w", err)
	}
	if schemaVal.Valid {
		et.Schema = json.RawMessage(schemaVal.String)
	}
	return et, nil
}

func (s *Store) DeleteEventType(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `DELETE FROM event_types WHERE id = $1`, id)
	return err
}

// ==================== Subscriptions ====================

func (s *Store) CreateSubscription(ctx context.Context, appID uuid.UUID, eventTypes []string, targetURL, description string) (*models.Subscription, error) {
	sub := &models.Subscription{}
	secret := "whsec_sub_" + uuid.New().String()
	eventTypesJSON, _ := json.Marshal(eventTypes)
	err := s.db.QueryRow(ctx,
		`INSERT INTO subscriptions (id, application_id, event_types, target_url, secret, description, enabled, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, true, $7, $7)
		 RETURNING id, application_id, event_types, target_url, secret, description, enabled, created_at, updated_at`,
		uuid.New(), appID, string(eventTypesJSON), targetURL, secret, description, time.Now(),
	).Scan(&sub.ID, &sub.ApplicationID, &eventTypesJSON, &sub.TargetURL, &sub.Secret, &sub.Description, &sub.Enabled, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}
	json.Unmarshal(eventTypesJSON, &sub.EventTypes)
	return sub, nil
}

func (s *Store) GetSubscription(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	sub := &models.Subscription{}
	var eventTypesJSON string
	err := s.db.QueryRow(ctx,
		`SELECT id, application_id, event_types, target_url, secret, description, enabled, created_at, updated_at
		 FROM subscriptions WHERE id = $1`, id,
	).Scan(&sub.ID, &sub.ApplicationID, &eventTypesJSON, &sub.TargetURL, &sub.Secret, &sub.Description, &sub.Enabled, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}
	json.Unmarshal([]byte(eventTypesJSON), &sub.EventTypes)
	return sub, nil
}

func (s *Store) ListSubscriptions(ctx context.Context, appID uuid.UUID, page, perPage int) ([]models.Subscription, int, error) {
	page, perPage = paginate(page, perPage)
	offset := (page - 1) * perPage

	var total int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM subscriptions WHERE application_id = $1`, appID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(ctx,
		`SELECT id, application_id, event_types, target_url, '', description, enabled, created_at, updated_at
		 FROM subscriptions WHERE application_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, appID, perPage, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		var eventTypesJSON string
		if err := rows.Scan(&sub.ID, &sub.ApplicationID, &eventTypesJSON, &sub.TargetURL, &sub.Secret, &sub.Description, &sub.Enabled, &sub.CreatedAt, &sub.UpdatedAt); err != nil {
			return nil, 0, err
		}
		json.Unmarshal([]byte(eventTypesJSON), &sub.EventTypes)
		subs = append(subs, sub)
	}
	return subs, total, nil
}

func (s *Store) ListSubscriptionsByEventType(ctx context.Context, appID uuid.UUID, eventType string) ([]models.Subscription, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, application_id, event_types, target_url, secret, description, enabled, created_at, updated_at
		 FROM subscriptions
		 WHERE application_id = $1 AND enabled = true AND event_types::jsonb ? $2`,
		appID, eventType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		var eventTypesJSON string
		if err := rows.Scan(&sub.ID, &sub.ApplicationID, &eventTypesJSON, &sub.TargetURL, &sub.Secret, &sub.Description, &sub.Enabled, &sub.CreatedAt, &sub.UpdatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(eventTypesJSON), &sub.EventTypes)
		subs = append(subs, sub)
	}
	return subs, nil
}

func (s *Store) UpdateSubscription(ctx context.Context, id uuid.UUID, req models.UpdateSubscriptionRequest) (*models.Subscription, error) {
	sub, err := s.GetSubscription(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.EventTypes != nil {
		sub.EventTypes = *req.EventTypes
	}
	if req.TargetURL != nil {
		sub.TargetURL = *req.TargetURL
	}
	if req.Description != nil {
		sub.Description = *req.Description
	}
	if req.Enabled != nil {
		sub.Enabled = *req.Enabled
	}

	eventTypesJSON, _ := json.Marshal(sub.EventTypes)
	_, err = s.db.Exec(ctx,
		`UPDATE subscriptions SET event_types = $2, target_url = $3, description = $4, enabled = $5, updated_at = $6 WHERE id = $1`,
		id, string(eventTypesJSON), sub.TargetURL, sub.Description, sub.Enabled, time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("update subscription: %w", err)
	}
	return sub, nil
}

func (s *Store) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `DELETE FROM subscriptions WHERE id = $1`, id)
	return err
}

// ==================== Events ====================

func (s *Store) CreateEvent(ctx context.Context, appID uuid.UUID, eventType string, payload, metadata json.RawMessage) (*models.Event, error) {
	e := &models.Event{}
	var metaVal pgtype.Text
	if metadata != nil {
		metaVal = pgtype.Text{String: string(metadata), Valid: true}
	}
	err := s.db.QueryRow(ctx,
		`INSERT INTO events (id, application_id, event_type, payload, metadata, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, application_id, event_type, payload, metadata, created_at`,
		uuid.New(), appID, eventType, string(payload), metaVal, time.Now(),
	).Scan(&e.ID, &e.ApplicationID, &e.EventType, &payload, &metaVal, &e.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}
	e.Payload = payload
	if metaVal.Valid {
		e.Metadata = json.RawMessage(metaVal.String)
	}
	return e, nil
}

func (s *Store) GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	e := &models.Event{}
	var payloadStr string
	var metaVal pgtype.Text
	err := s.db.QueryRow(ctx,
		`SELECT id, application_id, event_type, payload, metadata, created_at FROM events WHERE id = $1`, id,
	).Scan(&e.ID, &e.ApplicationID, &e.EventType, &payloadStr, &metaVal, &e.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}
	e.Payload = json.RawMessage(payloadStr)
	if metaVal.Valid {
		e.Metadata = json.RawMessage(metaVal.String)
	}
	return e, nil
}

func (s *Store) ListEvents(ctx context.Context, appID uuid.UUID, eventType string, page, perPage int) ([]models.Event, int, error) {
	page, perPage = paginate(page, perPage)
	offset := (page - 1) * perPage

	where := "WHERE application_id = $1"
	args := []interface{}{appID}
	argIdx := 2

	if eventType != "" {
		where += fmt.Sprintf(" AND event_type = $%d", argIdx)
		args = append(args, eventType)
		argIdx++
	}

	var total int
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM events `+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, perPage, offset)
	rows, err := s.db.Query(ctx,
		fmt.Sprintf(`SELECT id, application_id, event_type, payload, metadata, created_at
		 FROM events %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1),
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		var payloadStr string
		var metaVal pgtype.Text
		if err := rows.Scan(&e.ID, &e.ApplicationID, &e.EventType, &payloadStr, &metaVal, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		e.Payload = json.RawMessage(payloadStr)
		if metaVal.Valid {
			e.Metadata = json.RawMessage(metaVal.String)
		}
		events = append(events, e)
	}
	return events, total, nil
}

// ==================== Delivery Attempts ====================

func (s *Store) CreateDeliveryAttempt(ctx context.Context, eventID, subID uuid.UUID, status string, statusCode int, reqBody, respBody string, durationMs int64, attemptNum int) (*models.DeliveryAttempt, error) {
	da := &models.DeliveryAttempt{}
	err := s.db.QueryRow(ctx,
		`INSERT INTO delivery_attempts (id, event_id, subscription_id, status, status_code, request_body, response_body, duration_ms, attempt_number, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id, event_id, subscription_id, status, status_code, request_body, response_body, duration_ms, attempt_number, created_at`,
		uuid.New(), eventID, subID, status, statusCode, reqBody, respBody, durationMs, attemptNum, time.Now(),
	).Scan(&da.ID, &da.EventID, &da.SubscriptionID, &da.Status, &da.StatusCode, &da.RequestBody, &da.ResponseBody, &da.DurationMs, &da.AttemptNumber, &da.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create delivery attempt: %w", err)
	}
	return da, nil
}

func (s *Store) GetDeliveryAttempt(ctx context.Context, id uuid.UUID) (*models.DeliveryAttempt, error) {
	da := &models.DeliveryAttempt{}
	err := s.db.QueryRow(ctx,
		`SELECT id, event_id, subscription_id, status, status_code, request_body, response_body, duration_ms, attempt_number, created_at
		 FROM delivery_attempts WHERE id = $1`, id,
	).Scan(&da.ID, &da.EventID, &da.SubscriptionID, &da.Status, &da.StatusCode, &da.RequestBody, &da.ResponseBody, &da.DurationMs, &da.AttemptNumber, &da.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get delivery attempt: %w", err)
	}
	return da, nil
}

func (s *Store) ListDeliveryAttempts(ctx context.Context, appID uuid.UUID, status string, page, perPage int) ([]models.DeliveryAttempt, int, error) {
	page, perPage = paginate(page, perPage)
	offset := (page - 1) * perPage

	where := `FROM delivery_attempts da
			  JOIN events e ON da.event_id = e.id
			  WHERE e.application_id = $1`
	args := []interface{}{appID}
	argIdx := 2

	if status != "" {
		where += fmt.Sprintf(" AND da.status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	var total int
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) `+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, perPage, offset)
	rows, err := s.db.Query(ctx,
		fmt.Sprintf(`SELECT da.id, da.event_id, da.subscription_id, da.status, da.status_code, da.request_body, da.response_body, da.duration_ms, da.attempt_number, da.created_at
		 %s ORDER BY da.created_at DESC LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1),
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var attempts []models.DeliveryAttempt
	for rows.Next() {
		var da models.DeliveryAttempt
		if err := rows.Scan(&da.ID, &da.EventID, &da.SubscriptionID, &da.Status, &da.StatusCode, &da.RequestBody, &da.ResponseBody, &da.DurationMs, &da.AttemptNumber, &da.CreatedAt); err != nil {
			return nil, 0, err
		}
		attempts = append(attempts, da)
	}
	return attempts, total, nil
}

func (s *Store) ListDeliveryAttemptsByEvent(ctx context.Context, eventID uuid.UUID) ([]models.DeliveryAttempt, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, event_id, subscription_id, status, status_code, request_body, response_body, duration_ms, attempt_number, created_at
		 FROM delivery_attempts WHERE event_id = $1 ORDER BY attempt_number ASC`, eventID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []models.DeliveryAttempt
	for rows.Next() {
		var da models.DeliveryAttempt
		if err := rows.Scan(&da.ID, &da.EventID, &da.SubscriptionID, &da.Status, &da.StatusCode, &da.RequestBody, &da.ResponseBody, &da.DurationMs, &da.AttemptNumber, &da.CreatedAt); err != nil {
			return nil, err
		}
		attempts = append(attempts, da)
	}
	return attempts, nil
}

// ==================== Stats ====================

func (s *Store) GetAppStats(ctx context.Context, appID uuid.UUID) (*models.StatsResponse, error) {
	stats := &models.StatsResponse{}

	s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM events WHERE application_id = $1`, appID,
	).Scan(&stats.TotalEvents)

	s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM delivery_attempts da JOIN events e ON da.event_id = e.id WHERE e.application_id = $1`, appID,
	).Scan(&stats.TotalDeliveries)

	s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM subscriptions WHERE application_id = $1`, appID,
	).Scan(&stats.TotalSubscriptions)

	var successCount int64
	s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM delivery_attempts da JOIN events e ON da.event_id = e.id WHERE e.application_id = $1 AND da.status = 'success'`, appID,
	).Scan(&successCount)

	if stats.TotalDeliveries > 0 {
		stats.SuccessRate = float64(successCount) / float64(stats.TotalDeliveries) * 100
	}

	return stats, nil
}

func (s *Store) GetOrgStats(ctx context.Context, orgID uuid.UUID) (*models.StatsResponse, error) {
	stats := &models.StatsResponse{}

	s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM events e JOIN applications a ON e.application_id = a.id WHERE a.organization_id = $1`, orgID,
	).Scan(&stats.TotalEvents)

	s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM delivery_attempts da JOIN events e ON da.event_id = e.id JOIN applications a ON e.application_id = a.id WHERE a.organization_id = $1`, orgID,
	).Scan(&stats.TotalDeliveries)

	s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM subscriptions s JOIN applications a ON s.application_id = a.id WHERE a.organization_id = $1`, orgID,
	).Scan(&stats.TotalSubscriptions)

	s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM applications WHERE organization_id = $1`, orgID,
	).Scan(&stats.TotalApplications)

	var successCount int64
	s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM delivery_attempts da JOIN events e ON da.event_id = e.id JOIN applications a ON e.application_id = a.id WHERE a.organization_id = $1 AND da.status = 'success'`, orgID,
	).Scan(&successCount)

	if stats.TotalDeliveries > 0 {
		stats.SuccessRate = float64(successCount) / float64(stats.TotalDeliveries) * 100
	}

	return stats, nil
}

func (s *Store) GetChartData(ctx context.Context, appID uuid.UUID, days int) ([]models.ChartDataPoint, error) {
	rows, err := s.db.Query(ctx,
		`SELECT
			DATE(da.created_at) as date,
			COUNT(*) FILTER (WHERE da.status = 'success') as success,
			COUNT(*) FILTER (WHERE da.status = 'failed') as failed,
			COUNT(*) as total
		 FROM delivery_attempts da
		 JOIN events e ON da.event_id = e.id
		 WHERE e.application_id = $1 AND da.created_at > NOW() - INTERVAL '1 day' * $2
		 GROUP BY DATE(da.created_at)
		 ORDER BY date ASC`, appID, days,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []models.ChartDataPoint
	for rows.Next() {
		var p models.ChartDataPoint
		if err := rows.Scan(&p.Date, &p.Success, &p.Failed, &p.Total); err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	return points, nil
}

// ==================== Auth Helpers ====================

func (s *Store) GetUserOrganizations(ctx context.Context, userID uuid.UUID) ([]models.Organization, error) {
	rows, err := s.db.Query(ctx,
		`SELECT o.id, o.name, o.owner_id, o.created_at, o.updated_at
		 FROM organizations o
		 JOIN organization_members om ON o.id = om.organization_id
		 WHERE om.user_id = $1
		 ORDER BY o.created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []models.Organization
	for rows.Next() {
		var o models.Organization
		if err := rows.Scan(&o.ID, &o.Name, &o.OwnerID, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, o)
	}
	return orgs, nil
}

// Application by API key lookup
func (s *Store) GetApplicationByAPIKey(ctx context.Context, key string) (*models.Application, error) {
	a := &models.Application{}
	err := s.db.QueryRow(ctx,
		`SELECT a.id, a.organization_id, a.name, a.description, a.created_at, a.updated_at
		 FROM applications a
		 JOIN application_secrets s ON a.id = s.application_id
		 WHERE s.key = $1`, key,
	).Scan(&a.ID, &a.OrganizationID, &a.Name, &a.Description, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get app by api key: %w", err)
	}
	return a, nil
}

// ==================== Invitations ====================

type Invitation struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	Token          string    `json:"token"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

func (s *Store) CreateInvitation(ctx context.Context, orgID uuid.UUID, email, role string) (*Invitation, error) {
	inv := &Invitation{}
	err := s.db.QueryRow(ctx,
		`INSERT INTO invitations (id, organization_id, email, role, token, created_at, expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, organization_id, email, role, token, created_at, expires_at`,
		uuid.New(), orgID, email, role, uuid.New().String(), time.Now(), time.Now().Add(7*24*time.Hour),
	).Scan(&inv.ID, &inv.OrganizationID, &inv.Email, &inv.Role, &inv.Token, &inv.CreatedAt, &inv.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("create invitation: %w", err)
	}
	return inv, nil
}

func (s *Store) GetInvitationByToken(ctx context.Context, token string) (*Invitation, error) {
	inv := &Invitation{}
	err := s.db.QueryRow(ctx,
		`SELECT id, organization_id, email, role, token, created_at, expires_at
		 FROM invitations WHERE token = $1 AND expires_at > NOW()`, token,
	).Scan(&inv.ID, &inv.OrganizationID, &inv.Email, &inv.Role, &inv.Token, &inv.CreatedAt, &inv.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("get invitation: %w", err)
	}
	return inv, nil
}

func (s *Store) DeleteInvitation(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `DELETE FROM invitations WHERE id = $1`, id)
	return err
}

// ListMembersWithUserDetails returns members with user info
type MemberWithUser struct {
	models.OrganizationMember
	UserEmail string `json:"user_email"`
	UserName  string `json:"user_name"`
}

func (s *Store) ListMembersWithUserDetails(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]MemberWithUser, int, error) {
	page, perPage = paginate(page, perPage)
	offset := (page - 1) * perPage

	var total int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM organization_members WHERE organization_id = $1`, orgID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(ctx,
		`SELECT om.id, om.organization_id, om.user_id, om.role, om.created_at, u.email, u.name
		 FROM organization_members om
		 JOIN users u ON om.user_id = u.id
		 WHERE om.organization_id = $1
		 ORDER BY om.created_at ASC LIMIT $2 OFFSET $3`, orgID, perPage, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var members []MemberWithUser
	for rows.Next() {
		var m MemberWithUser
		if err := rows.Scan(&m.ID, &m.OrganizationID, &m.UserID, &m.Role, &m.CreatedAt, &m.UserEmail, &m.UserName); err != nil {
			return nil, 0, err
		}
		members = append(members, m)
	}
	return members, total, nil
}

// ListApplicationSecretsFiltered for listing with optional app filter
func (s *Store) ListApplicationSecretsFiltered(ctx context.Context, appID uuid.UUID, page, perPage int) ([]models.ApplicationSecret, int, error) {
	return s.ListApplicationSecrets(ctx, appID, page, perPage)
}

// ListDeliveryAttemptsBySubscription for retry logic
func (s *Store) ListDeliveryAttemptsBySubscription(ctx context.Context, eventID, subID uuid.UUID) ([]models.DeliveryAttempt, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, event_id, subscription_id, status, status_code, request_body, response_body, duration_ms, attempt_number, created_at
		 FROM delivery_attempts WHERE event_id = $1 AND subscription_id = $2 ORDER BY attempt_number ASC`, eventID, subID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []models.DeliveryAttempt
	for rows.Next() {
		var da models.DeliveryAttempt
		if err := rows.Scan(&da.ID, &da.EventID, &da.SubscriptionID, &da.Status, &da.StatusCode, &da.RequestBody, &da.ResponseBody, &da.DurationMs, &da.AttemptNumber, &da.CreatedAt); err != nil {
			return nil, err
		}
		attempts = append(attempts, da)
	}
	return attempts, nil
}

// CountAppEvents counts events for an application
func (s *Store) CountAppEvents(ctx context.Context, appID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM events WHERE application_id = $1`, appID).Scan(&count)
	return count, err
}

// CountOrgEvents counts events for an organization
func (s *Store) CountOrgEvents(ctx context.Context, orgID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM events e JOIN applications a ON e.application_id = a.id WHERE a.organization_id = $1`, orgID,
	).Scan(&count)
	return count, err
}

// CountAppDeliveries counts deliveries for an application
func (s *Store) CountAppDeliveries(ctx context.Context, appID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM delivery_attempts da JOIN events e ON da.event_id = e.id WHERE e.application_id = $1`, appID,
	).Scan(&count)
	return count, err
}

// CountAppSubscriptions counts subscriptions for an application
func (s *Store) CountAppSubscriptions(ctx context.Context, appID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM subscriptions WHERE application_id = $1`, appID).Scan(&count)
	return count, err
}

// GetAppSuccessRate calculates success rate for an application
func (s *Store) GetAppSuccessRate(ctx context.Context, appID uuid.UUID) (float64, error) {
	var total, success int64
	s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM delivery_attempts da JOIN events e ON da.event_id = e.id WHERE e.application_id = $1`, appID,
	).Scan(&total)
	s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM delivery_attempts da JOIN events e ON da.event_id = e.id WHERE e.application_id = $1 AND da.status = 'success'`, appID,
	).Scan(&success)
	if total == 0 {
		return 0, nil
	}
	return float64(success) / float64(total) * 100, nil
}

// Helper to join strings
func joinStrings(strs []string, sep string) string {
	return strings.Join(strs, sep)
}

// ==================== Verification Tokens ====================

func (s *Store) CreateVerificationToken(ctx context.Context, userID uuid.UUID, email string) (string, error) {
	token := "verify_" + uuid.New().String()
	_, err := s.db.Exec(ctx,
		`INSERT INTO verification_tokens (id, user_id, token, email, created_at, expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New(), userID, token, email, time.Now(), time.Now().Add(24*time.Hour),
	)
	if err != nil {
		return "", fmt.Errorf("create verification token: %w", err)
	}
	return token, nil
}

func (s *Store) GetVerificationToken(ctx context.Context, token string) (uuid.UUID, string, error) {
	var userID uuid.UUID
	var email string
	err := s.db.QueryRow(ctx,
		`SELECT user_id, email FROM verification_tokens WHERE token = $1 AND expires_at > NOW()`, token,
	).Scan(&userID, &email)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("get verification token: %w", err)
	}
	return userID, email, nil
}

func (s *Store) DeleteVerificationToken(ctx context.Context, token string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM verification_tokens WHERE token = $1`, token)
	return err
}

func (s *Store) MarkUserVerified(ctx context.Context, userID uuid.UUID) error {
	// We'll store verification status in a new column
	_, err := s.db.Exec(ctx,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified BOOLEAN NOT NULL DEFAULT false;
		 UPDATE users SET email_verified = true, updated_at = $2 WHERE id = $1`,
		userID, time.Now(),
	)
	return err
}

func (s *Store) IsUserVerified(ctx context.Context, userID uuid.UUID) (bool, error) {
	// Ensure column exists
	s.db.Exec(ctx, `ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified BOOLEAN NOT NULL DEFAULT false`)

	var verified bool
	err := s.db.QueryRow(ctx,
		`SELECT email_verified FROM users WHERE id = $1`, userID,
	).Scan(&verified)
	return verified, err
}
