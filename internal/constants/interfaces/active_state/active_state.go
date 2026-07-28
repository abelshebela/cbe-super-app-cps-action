package active_state

import "net/http"

type ActiveState interface {
	UpdateStatus(w http.ResponseWriter, r *http.Request)
}
