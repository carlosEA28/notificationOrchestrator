package events

import (
	"encoding/json"
	"time"
)

type UserPreferencesUpdated struct {
	UserID    string                 `json:"userId"`
	EventType string                 `json:"eventType"`
	Channel   string                 `json:"channel"`
	Enabled   bool                   `json:"enabled"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

func NewUserPreferencesUpdated(userID, eventType, channel string, enabled bool) *UserPreferencesUpdated {
	return &UserPreferencesUpdated{
		UserID:    userID,
		EventType: eventType,
		Channel:   channel,
		Enabled:   enabled,
		Timestamp: time.Now(),
	}
}

func (e *UserPreferencesUpdated) UnmarshalJSON(data []byte) error {
	type Alias UserPreferencesUpdated
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
