package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/carlosEA28/notificationOrchestrator/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLNotificationRepository struct {
	db *pgxpool.Pool
}

func NewSQLNotificationRepository(db *pgxpool.Pool) *SQLNotificationRepository {
	return &SQLNotificationRepository{db: db}
}

const (
	queryGetUserWithPreferences = `
		SELECT 
			u.email, u.phone, u."pushToken",
			up.enabled, up.channel
		FROM "User" u
		INNER JOIN "UserPreferences" up ON up."userId" = u.id
		WHERE u.id = $1 AND up."eventType" = $2
	`

	queryGetTemplateBySlug = `
		SELECT "id", "slug", "channel", "subject", "content"
		FROM "NotificationTemplate"
		WHERE "slug" = $1
	`

	queryGetTemplateByID = `
		SELECT "id", "slug", "channel", "subject", "content"
		FROM "NotificationTemplate"
		WHERE "id" = $1
	`

	queryUpdateStatus = `
		UPDATE "Notification"
		SET status = $2, "updatedAt" = NOW()
		WHERE id = $1
	`
)

type UserPreferenceResult struct {
	Email     string
	Phone     *string
	PushToken *string
	Enabled   bool
	Channel   string
}

func (r *SQLNotificationRepository) GetUserWithPreferences(ctx context.Context, userID string, eventType string) (*UserPreferenceResult, error) {
	rows, err := r.db.Query(ctx, queryGetUserWithPreferences, userID, eventType)
	if err != nil {
		return nil, fmt.Errorf("failed to query user with preferences: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var result UserPreferenceResult
		if err := rows.Scan(&result.Email, &result.Phone, &result.PushToken, &result.Enabled, &result.Channel); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		return &result, nil
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return nil, pgx.ErrNoRows
}

func (r *SQLNotificationRepository) GetTemplateBySlug(ctx context.Context, slug string) (*domain.NotificationTemplate, error) {
	var tpl domain.NotificationTemplate
	var subject *string

	err := r.db.QueryRow(ctx, queryGetTemplateBySlug, slug).Scan(
		&tpl.ID, &tpl.Slug, &tpl.Channel, &subject, &tpl.Content,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	tpl.Subject = subject
	return &tpl, nil
}

func (r *SQLNotificationRepository) GetTemplateByID(ctx context.Context, id string) (*domain.NotificationTemplate, error) {
	var tpl domain.NotificationTemplate
	var subject *string

	err := r.db.QueryRow(ctx, queryGetTemplateByID, id).Scan(
		&tpl.ID, &tpl.Slug, &tpl.Channel, &subject, &tpl.Content,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get template by id: %w", err)
	}

	tpl.Subject = subject
	return &tpl, nil
}

func (r *SQLNotificationRepository) UpdateStatus(ctx context.Context, id string, status domain.NotificationStatus, providerResponse string) error {
	var responseJSON map[string]interface{}
	if providerResponse != "" {
		if err := json.Unmarshal([]byte(providerResponse), &responseJSON); err != nil {
			return fmt.Errorf("failed to parse provider response: %w", err)
		}
	}

	_, err := r.db.Exec(ctx, queryUpdateStatus, id, string(status))
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}
