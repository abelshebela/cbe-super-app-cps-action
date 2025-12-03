package feedback

import "net/http"

type FeedbackAdapter interface {
	CreateFeedback(w http.ResponseWriter, r *http.Request)
	GetFeedbacks(w http.ResponseWriter, r *http.Request)
	GetFeedbackByID(w http.ResponseWriter, r *http.Request)
}
