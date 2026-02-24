package donation

import "net/http"

type DonationHandler interface {
	CreateDonation(w http.ResponseWriter, r *http.Request)
	UpdateDonation(w http.ResponseWriter, r *http.Request)
	FetchDonation(w http.ResponseWriter, r *http.Request)
	FetchDonationByID(w http.ResponseWriter, r *http.Request)

	UpdateDonationImage(w http.ResponseWriter, r *http.Request)
	DeleteDonationImage(w http.ResponseWriter, r *http.Request)
	AddDonationImage(w http.ResponseWriter, r *http.Request)
	EnableDonation(w http.ResponseWriter, r *http.Request)
	DisableDonation(w http.ResponseWriter, r *http.Request)
}
