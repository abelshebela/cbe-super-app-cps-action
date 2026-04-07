package notification

import "net/http"

type NotificationHandler interface {
	CreateNotification(w http.ResponseWriter, r *http.Request)
	UpdateNotification(w http.ResponseWriter, r *http.Request)
	DeleteNotification(w http.ResponseWriter, r *http.Request)
	EnableNotification(w http.ResponseWriter, r *http.Request)
	DisableNotification(w http.ResponseWriter, r *http.Request)

	FetchNotificationByID(w http.ResponseWriter, r *http.Request)
	FetchNotifications(w http.ResponseWriter, r *http.Request)
}
