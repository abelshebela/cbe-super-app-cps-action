package customerkyc

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"
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
	CustomerStatus       string               `json:"customer_status"`

	KYC             KycInfo    `bson:"kyc" json:"kyc"`
	CorePhoneNumber string     `json:"core_phone_number"`
	ComplyCube      ComplyCube `bson:"complycube" json:"complycube"`

	MoneyLaunderingFree *bool  `json:"money_laundering_free"`
	KYCRejectReason     string `bson:"kyc_reject_reason,omitempty" json:"kyc_reject_reason,omitempty"`
	TermsAndConditions  string `json:"terms_and_conditions,omitempty"`

	Reviewer *imodel.UserInfo  `json:"reviewer,omitempty"`
	Review   *imodel.KYCReview `json:"review"`

	LinkedAccount []imodel.LinkedAccounts `json:"linked_account,omitempty"`

	ServerTime string `json:"server-time,omitempty"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
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

type KycInfo struct {
	PhoneMismatch  bool   `bson:"phone_mismatch" json:"phone_mismatch"`
	BelowThreshold bool   `bson:"below_threshold" json:"below_threshold"`
	ContainsANDOR  bool   `bson:"contains_and_or" json:"contains_and_or"`
	KYCStatus      string `bson:"kyc_status" json:"kyc_status"`
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

type ComplyCube struct {
	DocumentType    string `bson:"document_type" json:"document_type"`
	IdentityOutcome string `bson:"identity_outcome" json:"identity_outcome"`
	IdentityStatus  string `bson:"identity_status" json:"identity_status"`
}
