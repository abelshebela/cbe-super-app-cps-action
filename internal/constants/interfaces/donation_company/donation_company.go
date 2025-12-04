package donation_company

import "net/http"

type DonationCompanyAdapter interface {
	CreateDonationCompany(w http.ResponseWriter, r *http.Request)
	UpdateDonationCompany(w http.ResponseWriter, r *http.Request)
	FetchDonationCompany(w http.ResponseWriter, r *http.Request)
	FetchDonationCompanyByID(w http.ResponseWriter, r *http.Request)
	AccountLookup(w http.ResponseWriter, r *http.Request)
	EnableDonationCompany(w http.ResponseWriter, r *http.Request)
	DisableDonationCompany(w http.ResponseWriter, r *http.Request)
}
