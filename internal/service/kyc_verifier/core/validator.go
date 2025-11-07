package kyc_verifier

import (
	"errors"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/dto/kyc_verifier"
)

func ValidateUpdate(req kyc_verifier.UpdateKYCRequest) error {
	return nil
}

func ValidateApprove(req kyc_verifier.ApproveKYCRequest) error {
	if !req.Approve && req.Reason == "" {
		return errors.New(localization.ErrorInvalidRequest.Code)
	}
	return nil
}
