package sitota

import "net/http"

type SitotaAdapter interface {
	GetAllSitotas(w http.ResponseWriter, r *http.Request)
	GetSitota(w http.ResponseWriter, r *http.Request)
}
