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

type CustomerKYC struct {
	ID         bson.ObjectID   `json:"id" bson:"_id,omitempty"`
	UserID     string          `json:"user_id" bson:"user_id,omitempty"`
	ClientID   string          `json:"client_id" bson:"client_id,omitempty"`
	KYCData    KYC             `json:"kyc_data" bson:"kyc_data"`
	ComplyCube *ComplyCubeData `json:"complycube,omitempty" bson:"complycube,omitempty"`
	// FaydaVerification    *FaydaVerification  `json:"fayda_verification,omitempty" bson:"fayda_verification,omitempty"`
	ReviewStatus         string              `json:"review_status" bson:"review_status,omitempty"`
	ReviewComments       string              `json:"review_comments" bson:"review_comments,omitempty"`
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

// KYC Review related models
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

type Address struct {
	Zone   string `json:"zone" bson:"zone"`
	Kebele string `json:"kebele" bson:"kebele"`
	Woreda string `json:"woreda" bson:"woreda"`
	Region string `json:"region" bson:"region"`
}

type KYC struct {
	Sub               string    `json:"sub" bson:"sub"`
	FullName          string    `json:"full_name" bson:"full_name"`
	FirstName         string    `json:"first_name" bson:"first_name,omitempty"`
	MiddleName        string    `json:"middle_name" bson:"middle_name,omitempty"`
	LastName          string    `json:"last_name" bson:"last_name,omitempty"`
	Email             string    `json:"email" bson:"email"`
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
	SubAccountType    string    `json:"sub_account_type" bson:"sub_account_type"`
	Currency          string    `json:"currency" bson:"currency"`
	SourceOfIncome    string    `json:"source_of_income" bson:"source_of_income"`
	MaritalStatus     string    `json:"marital_status" bson:"marital_status"`
	Occupation        string    `json:"occupation" bson:"occupation"`
	OriginID          string    `json:"origin_id" bson:"origin_id"`
	IsCitizen         bool      `json:"is_citizen" bson:"is_citizen"`
	USTIN             string    `json:"us_tin" bson:"us_tin"`
	Country           string    `json:"country" bson:"country"`
	MothersName       string    `json:"mothers_name" bson:"mothers_name"`
	Vendor            Vendor    `json:"vendor" bson:"vendor"`
	Address           Address   `json:"address" bson:"address"`
}

type CustomerBarUnBarReason struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string        `json:"user_id" bson:"user_id"`
	Reason    string        `json:"reason" bson:"reason"`
	IsBarred  bool          `json:"is_barred" bson:"is_barred"`
	CreatedBy string        `json:"created_by" bson:"created_by"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
}

type ComplyCubeData struct {
	DocumentID      string    `json:"document_id,omitempty" bson:"document_id,omitempty"`
	LivePhotoID     string    `json:"live_photo_id,omitempty" bson:"live_photo_id,omitempty"`
	LiveVideoID     string    `json:"live_video_id,omitempty" bson:"live_video_id,omitempty"`
	DocumentType    string    `json:"document_type,omitempty" bson:"document_type,omitempty"`
	IdentityCheckID string    `json:"identity_check_id,omitempty" bson:"identity_check_id,omitempty"`
	DocumentCheckID string    `json:"document_check_id,omitempty" bson:"document_check_id,omitempty"`
	AMLCheckID      string    `json:"aml_check_id,omitempty" bson:"aml_check_id,omitempty"`
	IdentityCheck   any       `json:"identity_check,omitempty" bson:"identity_check,omitempty"`
	DocumentCheck   any       `json:"document_check,omitempty" bson:"document_check,omitempty"`
	AMLCheck        any       `json:"aml_check,omitempty" bson:"aml_check,omitempty"`
	Document        any       `json:"document,omitempty" bson:"document,omitempty"`
	ExtractedData   any       `json:"extracted_data,omitempty" bson:"extracted_data,omitempty"`
	IdentityOutcome string    `json:"identity_outcome,omitempty" bson:"identity_outcome,omitempty"`
	DocumentOutcome string    `json:"document_outcome,omitempty" bson:"document_outcome,omitempty"`
	AMLOutcome      string    `json:"aml_outcome,omitempty" bson:"aml_outcome,omitempty"`
	IdentityStatus  string    `json:"identity_status,omitempty" bson:"identity_status,omitempty"`
	DocumentStatus  string    `json:"document_status,omitempty" bson:"document_status,omitempty"`
	AMLStatus       string    `json:"aml_status,omitempty" bson:"aml_status,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
