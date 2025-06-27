package inbound

import "net/http"

type Inbound interface {
	FetchServices(w http.ResponseWriter, r *http.Request)
	EnableDisableServicesMaker(w http.ResponseWriter, r *http.Request)
	EnableDisableServicesChecker(w http.ResponseWriter, r *http.Request)

	SearchAccountByCif(w http.ResponseWriter, r *http.Request)
	RemoveCifMaker(w http.ResponseWriter, r *http.Request)
	RemoveCifChecker(w http.ResponseWriter, r *http.Request)
}
