package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/carlosEA28/notificationOrchestrator/internal/domain"
	"github.com/carlosEA28/notificationOrchestrator/internal/events"
	"github.com/carlosEA28/notificationOrchestrator/internal/repository"
)

type NotificationProcessor struct {
	repo *repository.SQLNotificationRepository
}

func NewNotificationProcessor(repo *repository.SQLNotificationRepository) *NotificationProcessor {
	return &NotificationProcessor{repo: repo}
}

type DeliveryPayload struct {
	Channel       string
	RoutingKey    string
	Payload       map[string]interface{}
	TemplateSlug  string
	TemplateID    string
	CorrelationID string
}

func (p *NotificationProcessor) BuildPayload(ctx context.Context, event events.NotificationRequested) (*DeliveryPayload, error) {
	templateSlug := event.TemplateSlug
	var template *domain.NotificationTemplate
	var err error

	eventType := event.EventType
	if eventType == "" {
		return nil, fmt.Errorf("evento incompleto: userId=%q eventType=%q templateSlug=%q", event.UserID, eventType, templateSlug)
	}

	if templateSlug == "" && event.TemplateID != "" {
		templateCtx, cancelTemplate := context.WithTimeout(ctx, 3*time.Second)
		defer cancelTemplate()

		template, err = p.repo.GetTemplateByID(templateCtx, event.TemplateID)
		if err != nil {
			return nil, fmt.Errorf("erro ao obter template por id: %w", err)
		}
		templateSlug = template.Slug
	}

	if event.UserID == "" || templateSlug == "" {
		return nil, fmt.Errorf("evento incompleto: userId=%q eventType=%q templateSlug=%q", event.UserID, eventType, templateSlug)
	}

	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	userWithPreferences, err := p.repo.GetUserWithPreferences(dbCtx, event.UserID, eventType)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter usuário com preferências: %w", err)
	}

	if !userWithPreferences.Enabled {
		return nil, nil
	}

	if template == nil {
		templateCtx, cancelTemplate := context.WithTimeout(ctx, 3*time.Second)
		defer cancelTemplate()
		template, err = p.repo.GetTemplateBySlug(templateCtx, templateSlug)
		if err != nil {
			return nil, fmt.Errorf("erro ao obter template: %w", err)
		}
	}

	content := replaceTemplateVariables(template.Content, event.Payload)

	routingKey := ""
	switch userWithPreferences.Channel {
	case "email":
		routingKey = "notification.delivery.email"
	case "push":
		routingKey = "notification.delivery.push"
	case "sms":
		routingKey = "notification.delivery.sms"
	default:
		return nil, fmt.Errorf("canal não suportado: %q", userWithPreferences.Channel)
	}

	finalPayload := map[string]interface{}{
		"email":         userWithPreferences.Email,
		"phone":         userWithPreferences.Phone,
		"pushToken":     userWithPreferences.PushToken,
		"content":       content,
		"variables":     event.Payload,
		"correlationId": event.CorrelationID,
		"channel":       userWithPreferences.Channel,
		"eventType":     eventType,
		"templateId":    event.TemplateID,
		"templateSlug":  templateSlug,
	}

	return &DeliveryPayload{
		Channel:       userWithPreferences.Channel,
		RoutingKey:    routingKey,
		Payload:       finalPayload,
		TemplateSlug:  templateSlug,
		TemplateID:    event.TemplateID,
		CorrelationID: event.CorrelationID,
	}, nil
}

func replaceTemplateVariables(content string, variables map[string]interface{}) string {
	if content == "" {
		return ""
	}

	result := content
	for key, value := range variables {
		// fmt.Sprint garante que mesmo números ou booleanos virem string
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, fmt.Sprint(value))
	}
	return result
}
