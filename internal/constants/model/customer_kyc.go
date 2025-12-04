package model

import (
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CustomerKYC struct {
	ID                   bson.ObjectID       `json:"id" bson:"_id"`
	UserID               string              `json:"user_id" bson:"user_id"`
	KYCData              types.KYCData       `json:"kyc_data" bson:"kyc_data"`
	KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed" bson:"kyc_reject_reason_failed"`
	KYCStatus            constants.KYCStatus `json:"kyc_status" bson:"kyc_status"`
	KYCRejectReason      string              `json:"kyc_reject_reason" bson:"kyc_reject_reason"`
	KYCApproved          bool                `json:"kyc_approved" bson:"kyc_approved"`
	KYCActivityBy        any                 `json:"kyc_activity_by" bson:"kyc_activity_by"`
	KYCLevel             uint8               `json:"kyc_level" bson:"kyc_level"`
	Enabled              bool                `json:"enabled" bson:"enabled"`
	IsDeleted            bool                `json:"is_deleted" bson:"is_deleted"`
	CreatedAt            time.Time           `json:"created_at" bson:"created_at"`
	LastModifiedAt       time.Time           `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt            time.Time           `json:"deleted_at" bson:"deleted_at"`
}
