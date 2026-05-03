package dto

type SendNotificationWorkerDTO struct {
	Email         string                 `json:"email"`
	Phone         string                 `json:"phone"`
	content       string                 `json:"content"`
	variables     map[string]interface{} `json:"variables"`
	correlationId string                 `json:"correlationId"`
}
