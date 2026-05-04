package events

import (
	"encoding/json"
	"time"
)

type NotificationRequested struct {
	NotificationID string                 `json:"notificationId,omitempty"`
	ID             string                 `json:"id,omitempty"`
	CorrelationID  string                 `json:"correlationId"`
	UserID         string                 `json:"userId"`
	TemplateSlug   string                 `json:"templateSlug,omitempty"`
	TemplateID     string                 `json:"templateId"`
	EventType      string                 `json:"eventType"`
	Payload        map[string]interface{} `json:"payload"`
	Priority       int                    `json:"priority,omitempty"`
	CreatedAt      time.Time              `json:"createdAt"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

func NewNotificationRequested(notificationID, correlationID, userID, templateSlug, eventType string, payload map[string]interface{}) *NotificationRequested {
	return &NotificationRequested{
		NotificationID: notificationID,
		CorrelationID:  correlationID,
		UserID:         userID,
		TemplateSlug:   templateSlug,
		EventType:      eventType,
		Payload:        payload,
		CreatedAt:      time.Now(),
	}
}

func (e *NotificationRequested) UnmarshalJSON(data []byte) error {
	type Alias NotificationRequested
	aux := &struct {
		CreatedAt string `json:"createdAt"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339Nano, aux.CreatedAt); err == nil {
			e.CreatedAt = t
		}
	}
	if e.NotificationID == "" && e.ID != "" {
		e.NotificationID = e.ID
	}
	return nil
}
