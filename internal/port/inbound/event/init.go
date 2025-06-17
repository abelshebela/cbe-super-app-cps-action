package event_inbound

import "net/http"

type Inbound interface {
	MakerCreateEvent(w http.ResponseWriter, r *http.Request)
	CheckerEvent(w http.ResponseWriter, r *http.Request)
	FetchEventById(w http.ResponseWriter, r *http.Request)
	FetchEvent(w http.ResponseWriter, r *http.Request)
}
