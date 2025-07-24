package miniappmerchant

import "net/http"

type MiniAppMerchantInbound interface {
	CreateMiniAppMerchant(w http.ResponseWriter, r *http.Request)
	DeleteMiniAppMerchant(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
	GetMiniAppMerchant(w http.ResponseWriter, r *http.Request)
	GetAllMiniAppMerchant(w http.ResponseWriter, r *http.Request)
	UpdateMiniAppMerchant(w http.ResponseWriter, r *http.Request)
}
