package updatedbulkservice

import "net/http"

type BulkServiceHandler interface {
	GetAllBulkServices(w http.ResponseWriter, r *http.Request)
	EnableBulkService(w http.ResponseWriter, r *http.Request)
	DisableBulkService(w http.ResponseWriter, r *http.Request)
}
