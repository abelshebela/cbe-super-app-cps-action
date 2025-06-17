package inbound

import "net/http"

type PortalCardBound interface {
	GetAllPortalCard(w http.ResponseWriter, r *http.Request)
}
