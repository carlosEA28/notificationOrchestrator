package events

import (
	"encoding/json"
	"time"
)

type UserCreated struct {
	UserID    string                 `json:"userId"`
	Email     string                 `json:"email"`
	Phone     *string                `json:"phone,omitempty"`
	PushToken *string                `json:"pushToken,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

func NewUserCreated(userID, email string, phone, pushToken *string) *UserCreated {
	return &UserCreated{
		UserID:    userID,
		Email:     email,
		Phone:     phone,
		PushToken: pushToken,
		Timestamp: time.Now(),
	}
}

func (e *UserCreated) UnmarshalJSON(data []byte) error {
	type Alias UserCreated
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
