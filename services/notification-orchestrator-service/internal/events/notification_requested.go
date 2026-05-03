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
	TemplateSlug   string                 `json:"templateSlug"`
	TemplateID     string                 `json:"templateId,omitempty"`
	EventType      string                 `json:"eventType"`
	Payload        map[string]interface{} `json:"payload"`
	Priority       int                    `json:"priority,omitempty"`
	Timestamp      time.Time              `json:"timestamp"`
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
		Timestamp:      time.Now(),
	}
}

func (e *NotificationRequested) UnmarshalJSON(data []byte) error {
	type Alias NotificationRequested
	aux := &struct {
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, aux.Timestamp); err == nil {
			e.Timestamp = t
		}
	}
	if e.NotificationID == "" && e.ID != "" {
		e.NotificationID = e.ID
	}
	return nil
}
