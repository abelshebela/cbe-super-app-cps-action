package cbesuperappcpsaction

import (
	"time"
)

type ReviewStatus string
type ReviewDecision string

const (
	DecisionApproved ReviewDecision = "APPROVED"
	DecisionRejected ReviewDecision = "REJECTED"
)

const (
	ReviewPending   ReviewStatus = "PENDING"
	ReviewStarted   ReviewStatus = "STARTED"
	ReviewCompleted ReviewStatus = "COMPLETED"
	ReviewCancelled ReviewStatus = "CANCELLED"
)

type KYCResponse struct {
	ID             string `json:"id"`
	CustomerNumber string `json:"customer_number"`

	AccountType string `json:"account_type"`
	Currency    string `json:"currency"`

	CustomerStatus string   `json:"customer_status"`
	AccountNumbers []string `json:"account_numbers"`

	PersonalInformation  *PersonalInformation  `json:"personal_information"`
	ResidentialAddress   *ResidentialAddress   `json:"residential_address"`
	FinancialInformation *FinancialInformation `json:"financial_information"`
	CapturedDocuments    *CapturedDocuments    `json:"captured_documents"`

	Review *KYCReview `json:"review"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PersonalInformation struct {
	FullName      string `json:"full_name"`
	MotherName    string `json:"mother_name"`
	PhoneNumber   string `json:"phone_number"`
	Email         string `json:"email"`
	Gender        string `json:"gender"`
	MaritalStatus string `json:"marital_status"`
	Nationality   string `json:"nationality"`
	DateOfBirth   string `json:"date_of_birth"`
	IsCitizen     bool   `json:"is_citizen"`
	OriginID      string `json:"origin_id"`
	USTIN         string `json:"us_tin"`
}

type CapturedDocuments struct {
	Photo         string `json:"photo"`
	LivenessVideo string `json:"liveness_video,omitempty"`
	IDCardFront   string `json:"id_card_front,omitempty"`
	IDCardBack    string `json:"id_card_back,omitempty"`
}

type FinancialInformation struct {
	EmploymentStatus     string `json:"employment_status"`
	Occupation           string `json:"occupation"`
	SourceOfIncome       string `json:"source_of_income"`
	AverageMonthlyIncome string `json:"average_monthly_income"`
}

type ResidentialAddress struct {
	Country     string `json:"country"`
	Region      string `json:"region"`
	Zone        string `json:"zone,omitempty"`
	City        string `json:"city,omitempty"`
	SubCity     string `json:"sub_city,omitempty"`
	Wereda      string `json:"wereda,omitempty"`
	Kebele      string `json:"kebele,omitempty"`
	HouseNumber string `json:"house_number,omitempty"`
}

type KYCReview struct {
	Assignments     []KYCReviewAssignment `json:"assignments"`
	PickCount       int                   `json:"pick_count"`
	Status          ReviewStatus          `json:"status"`
	Decision        *KYCReviewDecision    `json:"decision"`
	ReviewStartedAt *time.Time            `json:"review_started_at"`
	ReviewExpiresAt *time.Time            `json:"review_expires_at"`
}

type KYCReviewDecision struct {
	Outcome         ReviewDecision `json:"outcome"`
	RejectionReason string         `json:"rejection_reason"`
	ReviewedAt      *time.Time     `json:"reviewed_at"`
	ReviewedBy      *UserInfo      `json:"reviewed_by"`
}

type KYCReviewAssignment struct {
	PickedAt   *time.Time `json:"picked_at"`
	PickedBy   *UserInfo  `json:"picked_by"`
	PickReason string     `json:"pick_reason"`
}

type UserInfo struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
