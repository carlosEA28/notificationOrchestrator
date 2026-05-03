package domain

import (
	"time"
)

type NotificationStatus string

const (
	StatusPending   NotificationStatus = "PENDING"
	StatusSent      NotificationStatus = "SENT"
	StatusDelivered NotificationStatus = "DELIVERED"
	StatusFailed    NotificationStatus = "FAILED"
)

type Notification struct {
	ID            string                 `json:"id"`
	UserID        string                 `json:"userId"`
	TemplateID    string                 `json:"templateId"`
	Status        NotificationStatus     `json:"status"`
	Payload       map[string]interface{} `json:"payload"`
	CorrelationID string                 `json:"correlationId"`
	Priority      int                    `json:"priority"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
}

type NotificationLog struct {
	ID               string                 `json:"id"`
	NotificationID   string                 `json:"notificationId"`
	Status           NotificationStatus     `json:"status"`
	Channel          string                 `json:"channel"`
	ProviderResponse map[string]interface{} `json:"providerResponse,omitempty"`
	ErrorMessage     *string                `json:"errorMessage,omitempty"`
	CreatedAt        time.Time              `json:"createdAt"`
}
