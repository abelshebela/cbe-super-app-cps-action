package utility

import "net/http"

type UtilityInbound interface {
	GetByUniqueToken(w http.ResponseWriter, r *http.Request)
}
