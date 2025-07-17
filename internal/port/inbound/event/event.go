package event_inbound

import "net/http"

type EventHandler interface {
	MakerCreateEvent(w http.ResponseWriter, r *http.Request)
	CheckerEvent(w http.ResponseWriter, r *http.Request)
	FetchEventByID(w http.ResponseWriter, r *http.Request)
	FetchEvents(w http.ResponseWriter, r *http.Request)
}
