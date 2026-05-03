package service

import (
	"context"
	"fmt"

	"github.com/carlosEA28/notificationOrchestrator/internal/events"
	"github.com/carlosEA28/notificationOrchestrator/internal/repository"
	"github.com/carlosEA28/notificationOrchestrator/internal/domain"
)

type NotificationProcessor struct {
	repo *repository.SQLNotificationRepository
}

func NewNotificationProcessor(repo *repository.SQLNotificationRepository) *NotificationProcessor {
	return &NotificationProcessor{repo: repo}
}

func (p *NotificationProcessor) BuildPayload(ctx context.Context, event events.NotificationRequested) (map[string]interface{}, error) {
	templateSlug := event.TemplateSlug
	var template *domain.NotificationTemplate
	var err error

	if templateSlug == "" && event.TemplateID != "" {
		template, err = p.repo.GetTemplateByID(ctx, event.TemplateID)
		if err != nil {
			return nil, fmt.Errorf("erro ao obter template por id: %w", err)
		}
		templateSlug = template.Slug
	}

	eventType := event.EventType
	if eventType == "" {
		eventType = templateSlug
	}

	if event.UserID == "" || eventType == "" || templateSlug == "" {
		return nil, fmt.Errorf("evento incompleto: userId=%q eventType=%q templateSlug=%q", event.UserID, eventType, templateSlug)
	}

	userWithPreferences, err := p.repo.GetUserWithPreferences(ctx, event.UserID, eventType)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter usuário com preferências: %w", err)
	}

	if len(userWithPreferences.Preferences) == 0 {
		return nil, fmt.Errorf("usuário sem preferências para eventType=%q", eventType)
	}

	if !userWithPreferences.Preferences[0].Enabled {
		return nil, fmt.Errorf("preferência desabilitada para eventType=%q", eventType)
	}

	if template == nil {
		template, err = p.repo.GetTemplateBySlug(ctx, templateSlug)
		if err != nil {
			return nil, fmt.Errorf("erro ao obter template: %w", err)
		}
	}

	finalPayload := map[string]interface{}{
		"email":         userWithPreferences.User.Email,
		"phone":         userWithPreferences.User.Phone,
		"content":       template.Content,
		"variables":     event.Payload,
		"correlationId": event.CorrelationID,
	}

	return finalPayload, nil
}
