package portal_card

import "net/http"

type PortalCardAdapter interface {
	GetAllPortalCard(w http.ResponseWriter, r *http.Request)
}
