package inbound

import "net/http"

type Feedback interface {
	GetFeedbacks(w http.ResponseWriter, r *http.Request)
	GetFeedbackByID(w http.ResponseWriter, r *http.Request)
}


