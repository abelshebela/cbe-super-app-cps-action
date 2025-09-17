package donation_company

import "net/http"

type DonationCompanyAdapter interface {
	CreateDonationCompany(w http.ResponseWriter, r *http.Request)
	UpdateDonationCompany(w http.ResponseWriter, r *http.Request)
	FetchDonationCompany(w http.ResponseWriter, r *http.Request)
	FetchDonationCompanyByID(w http.ResponseWriter, r *http.Request)
}
