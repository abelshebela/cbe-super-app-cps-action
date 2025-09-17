package event_inbound

import "net/http"

type EventHandler interface {
	CreateEvent(w http.ResponseWriter, r *http.Request)
	UpdateEvent(w http.ResponseWriter, r *http.Request)
	DeleteEvent(w http.ResponseWriter, r *http.Request)
	EnableEvent(w http.ResponseWriter, r *http.Request)
	DisableEvent(w http.ResponseWriter, r *http.Request)

	FetchEventByID(w http.ResponseWriter, r *http.Request)
	FetchEvents(w http.ResponseWriter, r *http.Request)
}
