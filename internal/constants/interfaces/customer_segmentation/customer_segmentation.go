package customersegmentation

import "net/http"

type CustomerSegmentation interface {
	CreateCustomerSegmentation(w http.ResponseWriter, r *http.Request)
}
