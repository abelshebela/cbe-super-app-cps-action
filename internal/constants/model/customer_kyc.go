package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type (
	KYCStatus string
	Vendor    string
)

const (
	KYCStatusPending     KYCStatus = "PENDING"
	KYCStatusInReview    KYCStatus = "IN_REVIEW"
	KYCStatusApproved    KYCStatus = "APPROVED"
	KYCStatusRejected    KYCStatus = "REJECTED"
	KYCStatusExpired     KYCStatus = "EXPIRED"
	KYCStatusTransferred KYCStatus = "TRANSFERRED"

	Fayda    Vendor = "FAYDA"
	Verigram Vendor = "Verigram"
)

type KYC struct {
	Sub               string    `json:"sub" bson:"sub"`
	FullName          string    `json:"full_name" bson:"full_name"`
	PhoneNumber       string    `json:"phone_number" bson:"phone_number"`
	Gender            string    `json:"gender" bson:"gender"`
	Picture           string    `json:"picture" bson:"picture"`
	SelfiePhoto       string    `json:"selfie_photo" bson:"selfie_photo"`
	Nationality       string    `json:"nationality" bson:"nationality"`
	BirthDate         time.Time `json:"birth_date" bson:"birth_date"`
	DocumentFront     string    `json:"document_front" bson:"document_front"`
	DocumentBack      string    `json:"document_back" bson:"document_back"`
	EmployementStatus string    `json:"employement_status" bson:"employement_status"`
	MonthlyIncome     string    `json:"monthly_income" bson:"monthly_income"`
	AccountType       string    `json:"account_type" bson:"account_type"`
	Country           string    `json:"country" bson:"country"`
	MothersName       string    `json:"mothers_name" bson:"mothers_name"`
	Vendor            Vendor    `json:"vendor" bson:"vendor"`
	Address           Address   `json:"address" bson:"address"`
}

type Address struct {
	Zone   string `json:"zone" bson:"zone"`
	Kebele string `json:"kebele" bson:"kebele"`
	Woreda string `json:"woreda" bson:"woreda"`
	Region string `json:"region" bson:"region"`
}

type CustomerKYC struct {
	ID                   bson.ObjectID       `json:"id" bson:"_id,omitempty"`
	UserID               string              `json:"user_id" bson:"user_id,omitempty"`
	KYCData              KYC                 `json:"kyc_data" bson:"kyc_data"`
	KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed" bson:"kyc_reject_reason_failed"`
	KYCStatus            KYCStatus           `json:"kyc_status" bson:"kyc_status"`
	KYCRejectReason      string              `json:"kyc_reject_reason" bson:"kyc_reject_reason"`
	KYCApproved          bool                `json:"kyc_approved" bson:"kyc_approved"`
	KYCActivityBy        any                 `json:"kyc_activity_by" bson:"kyc_activity_by"`
	Enabled              bool                `json:"enabled" bson:"enabled"`
	IsDeleted            bool                `json:"is_deleted" bson:"is_deleted"`
	CreatedAt            time.Time           `json:"created_at" bson:"created_at"`
	LastModifiedAt       time.Time           `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt            time.Time           `json:"deleted_at" bson:"deleted_at"`
}

type UserInfo struct {
	ID          bson.ObjectID `json:"id" bson:"id"`
	UserCode    string        `json:"user_code" bson:"user_code"`
	FullName    string        `json:"full_name" bson:"full_name"`
	Email       string        `json:"email" bson:"email"`
	Department  string        `json:"department" bson:"department"`
	PhoneNumber string        `json:"phone_number" bson:"phone_number"`
}

type StartedKycReview struct {
	ID       bson.ObjectID `json:"id" bson:"_id,omitempty"`
	KycID    bson.ObjectID `json:"kyc_id" bson:"kyc_id"`
	Reviewer UserInfo      `json:"reviewer" bson:"reviewer"`

	ReviewStatus string `json:"review_status" bson:"review_status"`

	StartedAt time.Time `json:"started_at" bson:"started_at"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`

	PickedAt   *time.Time `json:"picked_at,omitempty" bson:"picked_at,omitempty"`
	PickedBy   *UserInfo  `json:"picked_by,omitempty" bson:"picked_by,omitempty"`
	PickReason string     `json:"pick_reason,omitempty" bson:"pick_reason,omitempty"`
	PickCount  int        `json:"pick_count" bson:"pick_count"`

	IsActive  bool `json:"is_active" bson:"is_active"`
	IsDeleted bool `json:"is_deleted" bson:"is_deleted"`

	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at" bson:"last_modified_at"`
}
