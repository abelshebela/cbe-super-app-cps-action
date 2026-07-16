package customerkyc

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CustomerKYCResponse struct {
	ID             string `json:"id"`
	SuperAppUserID string `json:"super_app_user_id,omitempty"`
	Sub            string `json:"sub,omitempty"`
	CustomerNumber string `bson:"customer_number" json:"customer_number"`
	AccountType    string `json:"account_type,omitempty"`
	SubAccountType string `json:"sub_account_type,omitempty"`
	Currency       string `json:"currency,omitempty"`
	Vendor         string `json:"vendor,omitempty"`

	PersonalInformation  PersonalInformation  `json:"personal_information"`
	ResidentialAddress   ResidentialAddress   `json:"residential_address"`
	FinancialInformation FinancialInformation `json:"financial_information"`
	CapturedDocuments    CapturedDocuments    `json:"captured_documents"`
	AccountNumbers       []string             `json:"account_numbers,omitempty"`
	CustomerStatus       string               `json:"customer_status,omitempty"`
	KYCStatus            string               `json:"kyc_status"`
	Review               KycReview            `bson:"kyc_review" json:"kyc_review"`
	MoneyLaunderingFree  *bool                `json:"money_laundering_free"`
	TermsAndConditions   string               `json:"terms_and_conditions,omitempty"`
	CreatedAt            string               `json:"created_at"`
	UpdatedAt            string               `json:"updated_at"`
}

type UserInfo struct {
	ID          bson.ObjectID `json:"id" bson:"id"`
	UserCode    string        `json:"user_code" bson:"user_code"`
	FullName    string        `json:"full_name" bson:"full_name"`
	Email       string        `json:"email" bson:"email"`
	Department  string        `json:"department" bson:"department"`
	PhoneNumber string        `json:"phone_number" bson:"phone_number"`
}

type KycReview struct {
	Reviewer     UserInfo   `json:"reviewer" bson:"reviewer"`
	ReviewStatus string     `json:"review_status" bson:"review_status"`
	StartedAt    time.Time  `json:"started_at" bson:"started_at"`
	ExpiresAt    time.Time  `json:"expires_at" bson:"expires_at"`
	PickedAt     *time.Time `json:"picked_at,omitempty" bson:"picked_at,omitempty"`
	PickedBy     *UserInfo  `json:"picked_by,omitempty" bson:"picked_by,omitempty"`
	PickReason   string     `json:"pick_reason,omitempty" bson:"pick_reason,omitempty"`
	PickCount    int        `json:"pick_count" bson:"pick_count"`
}

type PersonalInformation struct {
	FullName       string   `json:"full_name" bson:"full_name"`
	MotherName     string   `json:"mother_name" validate:"required"`
	PhoneNumber    string   `json:"phone_number" validate:"required"`
	Email          string   `json:"email,omitempty"`
	Gender         string   `json:"gender" validate:"required,oneof=MALE FEMALE"`
	MaritalStatus  string   `json:"marital_status,omitempty"`
	Nationality    string   `json:"nationality,omitempty"`
	DateOfBirth    string   `json:"date_of_birth" validate:"required"`
	IsCitizen      bool     `json:"is_citizen"`
	OriginID       string   `json:"origin_id,omitempty"`
	USTIN          string   `json:"us_tin,omitempty"`
	AccountNumbers []string `bson:"account_numbers,omitempty" json:"account_numbers,omitempty"`
}

type CapturedDocuments struct {
	Photo         string `json:"photo"`
	LivenessVideo string `json:"liveness_video,omitempty"`
	IDCardFront   string `json:"id_card_front"`
	IDCardBack    string `json:"id_card_back"`
}

type FinancialInformation struct {
	EmploymentStatus     string `json:"employment_status"`
	Occupation           string `json:"occupation"`
	SourceOfIncome       string `json:"source_of_income,omitempty"`
	AverageMonthlyIncome string `json:"average_monthly_income"`
}

type ResidentialAddress struct {
	Country     string `json:"country" validate:"required"`
	Region      string `json:"region" validate:"required"`
	Zone        string `json:"zone,omitempty"`
	City        string `json:"city,omitempty"`
	SubCity     string `json:"sub_city,omitempty"`
	Wereda      string `json:"wereda,omitempty"`
	Kebele      string `json:"kebele,omitempty"`
	HouseNumber string `json:"house_number,omitempty"`
}
