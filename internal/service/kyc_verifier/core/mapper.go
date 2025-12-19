package core

import (
	"cbe-super-app-cps-action/internal/constants/dto/kyc_verifier"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_constant "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
)

func MapUpdateToModel(prev *model.CustomerKYC, req kyc_verifier.UpdateKYCRequest) *model.CustomerKYC {
	m := *prev
	if req.KYCStatus != "" {
		m.KYCStatus = shared_constant.KYCStatus(req.KYCStatus)
	}
	if req.KYCRejectReason != "" {
		m.KYCRejectReason = req.KYCRejectReason
	}
	if req.KYCRejectReasonField != nil {
		m.KYCRejectReasonField = req.KYCRejectReasonField
	}
	if req.KYCApproved != nil {
		m.KYCApproved = *req.KYCApproved
	}
	if req.KYCActivityBy != nil {
		m.KYCActivityBy = req.KYCActivityBy
	}
	if req.KYCLevel != nil {
		m.KYCLevel = *req.KYCLevel
	}
	return &m
}
