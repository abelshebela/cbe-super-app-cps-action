package model

import (
	"time"

	"cbe-super-app-cps-action/internal/constants"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CustomerKYC struct {
	ID                   bson.ObjectID       `json:"id" bson:"_id,omitempty"`
	UserID               bson.ObjectID       `json:"user_id" bson:"user_id,omitempty"`
	KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed" bson:"kyc_reject_reason_failed,omitempty"`
	KYCStatus            constants.KYCStatus `json:"kyc_status" bson:"kyc_status,omitempty"`
	KYCRejectReason      string              `json:"kyc_reject_reason" bson:"kyc_reject_reason,omitempty"`
	KYCApproved          bool                `json:"kyc_approved" bson:"kyc_approved,omitempty"`
	KYCActivityBy        any                 `json:"kyc_activity_by" bson:"kyc_activity_by,omitempty"`
	KYCLevel             uint8               `json:"kyc_level" bson:"kyc_level,omitempty"`
	Enabled              bool                `json:"enabled" bson:"enabled"`
	IsDeleted            bool                `json:"is_deleted" bson:"is_deleted,omitempty"`
	CreatedAt            time.Time           `json:"created_at" bson:"created_at,omitempty"`
	LastModifiedAt       time.Time           `json:"last_modified_at" bson:"last_modified_at,omitempty"`
	DeletedAt            time.Time           `json:"deleted_at" bson:"deleted_at,omitempty"`
}
