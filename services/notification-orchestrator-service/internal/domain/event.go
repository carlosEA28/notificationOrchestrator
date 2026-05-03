package domain

type NotificationEvent struct {
	NotificationID string                 `json:"notificationId"`
	CorrelationID  string                 `json:"correlationId"`
	UserID         string                 `json:"userId"`
	TemplateSlug   string                 `json:"templateSlug"`
	EventType      string                 `json:"eventType"`
	Payload        map[string]interface{} `json:"payload"`
}
