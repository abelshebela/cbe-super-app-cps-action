package logistics_merchant_adaptor

import "net/http"

type LogisticMerchantInboundAdaptor interface {
	CreateLogisticMerchant(w http.ResponseWriter, r *http.Request)
	UpdateLogisticMerchant(w http.ResponseWriter, r *http.Request)
	EnableLogisticMerchant(w http.ResponseWriter, r *http.Request)
	DisableLogisticMerchant(w http.ResponseWriter, r *http.Request)
	DeleteLogisticMerchant(w http.ResponseWriter, r *http.Request)
	GetLogisticMerchantByID(w http.ResponseWriter, r *http.Request)
	GetLogisticMerchants(w http.ResponseWriter, r *http.Request)
}
