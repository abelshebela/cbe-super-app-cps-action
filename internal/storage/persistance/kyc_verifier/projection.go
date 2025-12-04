package kyc_verifier

import (
	"cbe-super-app-cps-action/internal/constants/dto/kyc_verifier"
	"cbe-super-app-cps-action/internal/constants/model"
)

func MapToResponse(m *model.CustomerKYC) *kyc_verifier.KYCVerifierResponse {
	return &kyc_verifier.KYCVerifierResponse{
		ID:                   m.ID,
		UserFullName:         "",
		UserPhoneNumber:      "",
		UserCustomerNumber:   "",
		UserCode:             "",
		KYCStatus:            string(m.KYCStatus),
		KYCRejectReason:      m.KYCRejectReason,
		KYCRejectReasonField: m.KYCRejectReasonField,
		KYCApproved:          m.KYCApproved,
		KYCLevel:             m.KYCLevel,
		CreatedAt:            m.CreatedAt,
	}
}

func MapToResponses(list []*model.CustomerKYC) []*kyc_verifier.KYCVerifierResponse {
	res := make([]*kyc_verifier.KYCVerifierResponse, len(list))
	for i, v := range list {
		res[i] = MapToResponse(v)
	}
	return res
}
