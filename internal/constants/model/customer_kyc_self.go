package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ReviewStatus string

const (
	KYCStatusPending   ReviewStatus = "PENDING"
	KYCStatusInReview  ReviewStatus = "IN_REVIEW"
	KYCStatusCancelled ReviewStatus = "CANCELLED"
	KYCStatusApproved  ReviewStatus = "APPROVED"
	KYCStatusRejected  ReviewStatus = "REJECTED"
)

type SelfActivationUser struct {
	ID             bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	Sub            string          `bson:"sub" json:"sub"`
	Name           string          `bson:"name" json:"name"`
	Email          string          `bson:"email" json:"email"`
	PhoneNumber    string          `bson:"phone_number" json:"phone_number"`
	Gender         string          `bson:"gender" json:"gender"`
	Picture        string          `bson:"picture" json:"picture"`
	Nationality    string          `bson:"nationality" json:"nationality"`
	BirthDate      string          `bson:"birth_date" json:"birth_date"`
	Address        CustomerAddress `bson:"address" json:"address"`
	Enabled        bool            `bson:"enabled" json:"enabled"`
	AccountNumbers []string        `bson:"account_numbers" json:"account_numbers"`
	ChosenAccounts []string        `bson:"chosen_accounts" json:"chosen_accounts"`
	CustomerNumber string          `bson:"customer_number" json:"customer_number"`
	ExpiryDate     string          `bson:"expiry_date" json:"expiry_date"`
	IssueDate      string          `bson:"issue_date" json:"issue_date"`
	SelfiePhoto    string          `json:"selfie_photo" bson:"selfie_photo"`

	KYC             KycInfo          `bson:"kyc" json:"kyc"`
	CorePhoneNumber string           `bson:"core_phone_number" json:"core_phone_number"`
	ComplyCube      ComplyCube       `bson:"complycube" json:"complycube"`
	LinkedAccount   []LinkedAccounts `json:"linked_accounts" bson:"linked_accounts"`

	Review *KYCReview `json:"review" bson:"review"`

	RegistrationDate time.Time `bson:"registration_time" json:"registration_time"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	LastModifiedAt   time.Time `bson:"last_modified_at" json:"last_modified_at"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at"`
}

type KYCReview struct {
	Assignments     []KYCReviewAssignment `json:"assignments,omitempty" bson:"assignments"`
	PickCount       int                   `json:"pick_count,omitempty" bson:"pick_count"`
	Status          ReviewStatus          `json:"status" bson:"status"`
	Decision        *KYCReviewDecision    `json:"decision,omitempty" bson:"decision"`
	ReviewStartedAt *time.Time            `json:"review_started_at" bson:"review_started_at"`
	ReviewExpiresAt *time.Time            `json:"review_expires_at" bson:"review_expires_at"`
}

type KYCReviewDecision struct {
	RejectionReason string     `json:"rejection_reason" bson:"rejection_reason"`
	ReviewedAt      *time.Time `json:"reviewed_at" bson:"reviewed_at"`
	Reviewer        *UserInfo  `json:"reviewer" bson:"reviewer"`
}

type KYCReviewAssignment struct {
	PickedAt   *time.Time `json:"picked_at" bson:"picked_at"`
	PickedBy   *UserInfo  `json:"picked_by" bson:"picked_by"`
	PickReason string     `json:"pick_reason" bson:"pick_reason"`
}

type LinkedAccounts struct {
	AccountHolderName  string `json:"account_holder_name" bson:"account_holder_name"`
	AccountNumber      string `json:"account_number" bson:"account_number"`
	AccountType        string `json:"account_type" bson:"account_type"`
	ProductCode        string `json:"product_code" bson:"product_code"`
	Currency           string `json:"currency" bson:"currency"`
	BranchName         string `json:"branch_name" bson:"branch_name"`
	BranchCode         string `json:"branch_code" bson:"branch_code"`
	InActive           string `json:"in_active" bson:"in_active"`
	PostingRestriction string `json:"posting_restriction" bson:"posting_restriction"`
	RestrictionType    string `json:"restriction_type" bson:"restriction_type"`
	AndorAccount       *bool  `json:"and_or_account" bson:"and_or_account"`
}

type CustomerAddress struct {
	Zone   string `bson:"zone" json:"zone"`
	Kebele string `bson:"kebele" json:"kebele"`
	Woreda string `bson:"woreda" json:"woreda"`
	Region string `bson:"region" json:"region"`
}

type KycInfo struct {
	PhoneMismatch  bool `bson:"phone_mismatch" json:"phone_mismatch"`
	BelowThreshold bool `bson:"below_threshold" json:"below_threshold"`
	ContainsANDOR  bool `bson:"contains_and_or" json:"contains_and_or"`
}

type ComplyCube struct {
	DocumentID      string `bson:"document_id" json:"document_id"`
	LiveVideoID     string `bson:"live_video_id" json:"live_video_id"`
	DocumentType    string `bson:"document_type" json:"document_type"`
	IdentityCheckID string `bson:"identity_check_id" json:"identity_check_id"`

	IdentityCheck IdentityCheck `bson:"identity_check" json:"identity_check"`

	IdentityOutcome string    `bson:"identity_outcome" json:"identity_outcome"`
	IdentityStatus  string    `bson:"identity_status" json:"identity_status"`
	UpdatedAt       time.Time `bson:"updated_at" json:"updated_at"`
}

type IdentityCheck struct {
	ID             string           `bson:"id" json:"id"`
	ClientID       string           `bson:"clientId" json:"clientId"`
	LiveVideoID    string           `bson:"liveVideoId" json:"liveVideoId"`
	DocumentID     string           `bson:"documentId" json:"documentId"`
	EntityName     string           `bson:"entityName" json:"entityName"`
	Type           string           `bson:"type" json:"type"`
	Status         string           `bson:"status" json:"status"`
	InitialOutcome string           `bson:"initialOutcome" json:"initialOutcome"`
	Result         IdentityResult   `bson:"result" json:"result"`
	Metadata       IdentityMetadata `bson:"metadata" json:"metadata"`
	CreatedAt      time.Time        `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time        `bson:"updatedAt" json:"updatedAt"`
}

type IdentityResult struct {
	Outcome   string            `bson:"outcome" json:"outcome"`
	Breakdown IdentityBreakdown `bson:"breakdown" json:"breakdown"`
}

type IdentityBreakdown struct {
	IntegrityAnalysis    IntegrityAnalysis    `bson:"integrityAnalysis" json:"integrityAnalysis"`
	FaceAnalysis         FaceAnalysis         `bson:"faceAnalysis" json:"faceAnalysis"`
	AuthenticityAnalysis AuthenticityAnalysis `bson:"authenticityAnalysis" json:"authenticityAnalysis"`
}

type IntegrityAnalysis struct {
	FaceDetection string `bson:"faceDetection" json:"faceDetection"`
}

type FaceAnalysis struct {
	FacialSimilarity       string                `bson:"facialSimilarity" json:"facialSimilarity"`
	PreviouslyEnrolledFace string                `bson:"previouslyEnrolledFace" json:"previouslyEnrolledFace"`
	Breakdown              FaceAnalysisBreakdown `bson:"breakdown" json:"breakdown"`
}

type FaceAnalysisBreakdown struct {
	FacialSimilarityScore int `bson:"facialSimilarityScore" json:"facialSimilarityScore"`
}

type AuthenticityAnalysis struct {
	SpoofedImageAnalysis            string                        `bson:"spoofedImageAnalysis" json:"spoofedImageAnalysis"`
	LivenessCheck                   string                        `bson:"livenessCheck" json:"livenessCheck"`
	LivenessVoiceChallengeAnalysis  string                        `bson:"livenessVoiceChallengeAnalysis" json:"livenessVoiceChallengeAnalysis"`
	LivenessActionChallengeAnalysis string                        `bson:"livenessActionChallengeAnalysis" json:"livenessActionChallengeAnalysis"`
	Breakdown                       AuthenticityAnalysisBreakdown `bson:"breakdown" json:"breakdown"`
}

type AuthenticityAnalysisBreakdown struct {
	LivenessCheckScore int `bson:"livenessCheckScore" json:"livenessCheckScore"`
}

type IdentityMetadata struct {
	LiveVideo LiveVideoMetadata `bson:"liveVideo" json:"liveVideo"`
}

type LiveVideoMetadata struct {
	Language string `bson:"language" json:"language"`
}

// Model from Export
type ExportSelfActivationRequest struct {
	CustomerName     string    `json:"customer_name" bson:"customer_name"`
	PhoneNumber      string    `json:"phone_number" bson:"phone_number"`
	Gender           string    `json:"gender" bson:"gender"`
	DateOfBirth      string    `json:"date_of_birth" bson:"date_of_birth"`
	Region           string    `json:"region" bson:"region"`
	RegistrationDate time.Time `json:"registration_date" bson:"registration_date"`
	RejectionReason  string    `json:"rejection_reason" bson:"rejection_reason"`
	CustomerStatus   string    `json:"customer_status" bson:"customer_status"`
	KYCStatus        string    `json:"kyc_status" bson:"kyc_status"`
}

// Action Log
type SelfActivationActionLog struct {
	ActionCode      string             `json:"action_code"`
	ActionID        string             `json:"action_id"`
	ActionRequest   string             `json:"action_request"`
	ActionType      string             `json:"action_type"`
	Maker           *MakerCheckerInfo  `json:"maker"`
	Checker         []MakerCheckerInfo `json:"checker"`
	Reason          string             `json:"reason"`
	Status          string             `json:"status"`
	ActionTakeAt    time.Time          `json:"action_take_at"`
	ActionUpdatedAt time.Time          `json:"action_updated_at"`
}

type MakerCheckerInfo struct {
	ID       string `json:"id" bson:"id"`
	FullName string `json:"full_name" bson:"full_name"`
}
