package domain

import (
	"time"
)

type NotificationTemplate struct {
	ID        string    `json:"id"`
	Slug      string    `json:"slug"`
	Channel   string    `json:"channel"`
	Subject   *string   `json:"subject,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
