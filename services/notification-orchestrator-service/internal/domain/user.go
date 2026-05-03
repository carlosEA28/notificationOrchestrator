package domain

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone,omitempty"`
	PushToken *string   `json:"pushToken,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserPreferences struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	EventType string `json:"eventType"`
	Channel   string `json:"channel"`
	Enabled   bool   `json:"enabled"`
}

type UserWithPreferences struct {
	User        User
	Preferences []UserPreferences
}
