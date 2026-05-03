package repository

import (
	"context"

	"github.com/carlosEA28/notificationOrchestrator/internal/domain"
)

type NotificationRepository interface {
	GetUserWithPreferences(ctx context.Context, userID string, eventType string) (*domain.UserWithPreferences, error)
	GetTemplateBySlug(ctx context.Context, slug string) (*domain.NotificationTemplate, error)
	GetTemplateByID(ctx context.Context, id string) (*domain.NotificationTemplate, error)
	UpdateStatus(ctx context.Context, id string, status domain.NotificationStatus, providerResponse string) error
}
