package kyc_verifier

import (
	"cbe-super-app-cps-action/internal/constants/dto/kyc_verifier"
	"net/http"
)

func ParseUpdateRequest(r *http.Request) (kyc_verifier.UpdateKYCRequest, error) {
	var req kyc_verifier.UpdateKYCRequest
	return req, nil
}

func ParseApproveRequest(r *http.Request) (kyc_verifier.ApproveKYCRequest, error) {
	var req kyc_verifier.ApproveKYCRequest
	return req, nil
}
