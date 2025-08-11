package donation_inbound

import "net/http"

type DonationHandler interface {
	CreateDonationCategory(w http.ResponseWriter, r *http.Request)
	UpdateDonationCategory(w http.ResponseWriter, r *http.Request)
	FetchDonationCategory(w http.ResponseWriter, r *http.Request)
	FetchDonationCategoryByID(w http.ResponseWriter, r *http.Request)

	// Donation Company methods
	CreateDonationCompany(w http.ResponseWriter, r *http.Request)
	UpdateDonationCompany(w http.ResponseWriter, r *http.Request)
	FetchDonationCompany(w http.ResponseWriter, r *http.Request)
	FetchDonationCompanyByID(w http.ResponseWriter, r *http.Request)

	// Donation methods
	CreateDonation(w http.ResponseWriter, r *http.Request)
	UpdateDonation(w http.ResponseWriter, r *http.Request)
	FetchDonation(w http.ResponseWriter, r *http.Request)
	FetchDonationByID(w http.ResponseWriter, r *http.Request)
}
