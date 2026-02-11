package event_merchant_port

import "net/http"

type EventMerchantInboundAdaptor interface {
	CreateEventMerchant(w http.ResponseWriter, r *http.Request)
	UpdateEventMerchant(w http.ResponseWriter, r *http.Request)
	EnableEventMerchant(w http.ResponseWriter, r *http.Request)
	DisableEventMerchant(w http.ResponseWriter, r *http.Request)
	DeleteEventMerchant(w http.ResponseWriter, r *http.Request)
	GetEventMerchantByID(w http.ResponseWriter, r *http.Request)
	GetEventMerchants(w http.ResponseWriter, r *http.Request)
}
