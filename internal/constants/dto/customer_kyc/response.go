package customerkyc

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CustomerKYCResponse struct {
	ID                   string               `json:"id"`
	PersonalInformation  PersonalInformation  `json:"personal_information"`
	ResidentialAddress   ResidentialAddress   `json:"residential_address"`
	FinancialInformation FinancialInformation `json:"financial_information"`
	CapturedDocuments    CapturedDocuments    `json:"captured_documents"`
	CustomerStatus       string               `json:"customer_status"`
	KYCStatus            string               `json:"kyc_status"`
	MoneyLaunderingFree  *bool                `json:"money_laundering_free"`
	TermsAndConditions   string               `json:"terms_and_conditions,omitempty"`
	KYCReviewStartedAt   *time.Time           `json:"started_at,omitempty"`
	KYCReviewExpiresAt   *time.Time           `json:"expires_at,omitempty"`
	ReviewerID           *bson.ObjectID       `json:"reviewer_id,omitempty"`
	CreatedAt            string               `json:"created_at"`
	UpdatedAt            string               `json:"updated_at"`
}

type PersonalInformation struct {
	FirstName    string `json:"first_name" validate:"required"`
	MiddleName   string `json:"middle_name,omitempty"`
	LastName     string `json:"last_name" validate:"required"`
	MotherName   string `json:"mother_name" validate:"required"`
	PhoneNumber  string `json:"phone_number" validate:"required"`
	Gender       string `json:"gender" validate:"required,oneof=MALE FEMALE"`
	MaritalStaus string `json:"marital_status"`
	Nationality  string `json:"nationality"`
	DateOfBirth  string `json:"date_of_birth" validate:"required"`
}

type CapturedDocuments struct {
	Photo         string `json:"photo"`
	LivenessVideo string `json:"liveness_video"`
	IDCardFront   string `json:"id_card_front"`
	IDCardBack    string `json:"id_card_back"`
}

type FinancialInformation struct {
	EmploymentStatus     string `json:"employment_status"`
	Occupation           string `json:"occupation"`
	AverageMonthlyIncome string `json:"average_monthly_income"`
}

type ResidentialAddress struct {
	Country     string `json:"country" validate:"required"`
	Region      string `json:"region" validate:"required"`
	City        string `json:"city" validate:"required"`
	SubCity     string `json:"sub_city,omitempty"`
	Wereda      string `json:"wereda,omitempty"`
	Kebele      string `json:"kebele,omitempty"`
	HouseNumber string `json:"house_number,omitempty"`
}
