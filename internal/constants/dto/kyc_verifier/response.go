package kyc_verifier

import (
	"time"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type KYCVerifierResponse struct {
	ID                   bson.ObjectID       `json:"id" bson:"id"`
	UserFullName         string              `json:"user_full_name" bson:"user_full_name"`
	UserPhoneNumber      string              `json:"phone_number" bson:"phone_number"`
	UserCustomerNumber   string              `json:"customer_number" bson:"customer_number"`
	UserCode             string              `json:"user_code" bson:"user_code"`
	KYCStatus            string              `json:"kyc_status" bson:"kyc_status"`
	KYCRejectReason      string              `json:"kyc_reject_reason" bson:"kyc_reject_reason"`
	KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed" bson:"kyc_reject_reason_failed"`
	KYCApproved          bool                `json:"kyc_approved" bson:"kyc_approved"`
	KYCActivityBy        any                 `json:"kyc_activity_by" bson:"kyc_activity_by"`
	KYCLevel             uint8               `json:"kyc_level" bson:"kyc_level"`
	CreatedAt            time.Time           `json:"created_at" bson:"created_at"`
}

type KYCVerifierListResponse KYCVerifierResponse
