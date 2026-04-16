package ussd_merchant_interface

import "net/http"

type UssdMerchantInbound interface {
	CreateUssdMerchant(w http.ResponseWriter, r *http.Request)
	UpdateUssdMerchant(w http.ResponseWriter, r *http.Request)
	EnableUssdMerchant(w http.ResponseWriter, r *http.Request)
	DisableUssdMerchant(w http.ResponseWriter, r *http.Request)
	DeleteUssdMerchant(w http.ResponseWriter, r *http.Request)
	GetUssdMerchant(w http.ResponseWriter, r *http.Request)
	GetAllUssdMerchant(w http.ResponseWriter, r *http.Request)
	DelereUssdMerchant(w http.ResponseWriter,r *http.Request)
}
