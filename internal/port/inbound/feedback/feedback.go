package inbound

import "net/http"

type Feedback interface {
	CreateFeedback(w http.ResponseWriter, r *http.Request)
	GetFeedbacks(w http.ResponseWriter, r *http.Request)
	GetFeedbackByID(w http.ResponseWriter, r *http.Request)
}
