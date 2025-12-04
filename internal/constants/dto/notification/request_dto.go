package notification

type NotificationRequest struct {
	NotificationType string `json:"notification_type"`
	NotificationBody string `json:"notification_body"`
	IsPublic         bool   `json:"is_public"`
	For              string `json:"for"`
	CreatedBy        string `json:"created_by"`
	Title            string `json:"title"`
}
