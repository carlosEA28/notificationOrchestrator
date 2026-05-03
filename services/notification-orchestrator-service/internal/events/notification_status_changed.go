package events

import (
	"encoding/json"
	"time"
)

type NotificationStatusChanged struct {
	NotificationID   string                 `json:"notificationId"`
	CorrelationID    string                 `json:"correlationId"`
	PreviousStatus   string                 `json:"previousStatus"`
	CurrentStatus    string                 `json:"currentStatus"`
	Channel          string                 `json:"channel,omitempty"`
	Provider         string                 `json:"provider,omitempty"`
	ProviderResponse map[string]interface{} `json:"providerResponse,omitempty"`
	ErrorMessage     *string                `json:"errorMessage,omitempty"`
	Timestamp        time.Time              `json:"timestamp"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

func NewNotificationStatusChanged(notificationID, correlationID, previousStatus, currentStatus string) *NotificationStatusChanged {
	return &NotificationStatusChanged{
		NotificationID: notificationID,
		CorrelationID:  correlationID,
		PreviousStatus: previousStatus,
		CurrentStatus:  currentStatus,
		Timestamp:      time.Now(),
	}
}

func (e *NotificationStatusChanged) UnmarshalJSON(data []byte) error {
	type Alias NotificationStatusChanged
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
	return nil
}
