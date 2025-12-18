package customer

import "net/http"

type CustomerDetail interface {
	GetCustomerDetail(w http.ResponseWriter, r *http.Request)
	GetCustomerByID(w http.ResponseWriter, r *http.Request)
	GetBlockedCustomer(w http.ResponseWriter, r *http.Request)
	SetEnableCustomerSession(w http.ResponseWriter, r *http.Request)
	ApproveFaydaCustomer(w http.ResponseWriter, r *http.Request)
	EnableCustomer(w http.ResponseWriter, r *http.Request)
	DisableCustomer(w http.ResponseWriter, r *http.Request)
	GetLinkedAccount(w http.ResponseWriter, r *http.Request)
	SearchCustomerByCIForAccountNumber(w http.ResponseWriter, r *http.Request)
}
