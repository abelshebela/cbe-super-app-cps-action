package donation_inbound

import "net/http"
type  DonationHandler interface {
	CreateDonationCategory(w http.ResponseWriter, r *http.Request)
	UpdateDonationCategory(w http.ResponseWriter, r *http.Request)
	FetchDonationCategory(w http.ResponseWriter, r *http.Request)
	FetchDonationCategoryByID(w http.ResponseWriter, r *http.Request)
}