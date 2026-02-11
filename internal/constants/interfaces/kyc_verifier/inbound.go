package kyc_verifier

import "net/http"

type KYCVerifierAdapter interface {
	GetKYCList(w http.ResponseWriter, r *http.Request)
	GetKYCByID(w http.ResponseWriter, r *http.Request)
	UpdateKYC(w http.ResponseWriter, r *http.Request)
	ApproveKYC(w http.ResponseWriter, r *http.Request)
}
