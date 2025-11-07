package kyc_verifier

import "mime/multipart"

type UpdateKYCRequest struct {
	KYCStatus            string                 `json:"kyc_status"`
	KYCRejectReason      string                 `json:"kyc_reject_reason"`
	KYCRejectReasonField map[string]struct{}    `json:"kyc_reject_reason_failed"`
	KYCApproved          *bool                  `json:"kyc_approved"`
	KYCActivityBy        any                    `json:"kyc_activity_by"`
	KYCLevel             *uint8                 `json:"kyc_level"`
	Attachment           *multipart.FileHeader  `json:"attachment"`
}

type ApproveKYCRequest struct {
	Approve bool   `json:"approve"`
	Reason  string `json:"reason"`
}
