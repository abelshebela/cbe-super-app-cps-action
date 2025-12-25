package cpsroles

import "net/http"

type CPSActionAdapter interface {
	createCPSRole(w http.ResponseWriter, r http.Request)
}
