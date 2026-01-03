package customer

import "cbe-super-app-cps-action/internal/constants"

type CustomerEnableDTO struct {
	UserOTP string `json:"user_otp"`
}
type CustomerDisableDTO struct {
	IsTemporary   *bool  `json:"is_temporary"`
	DisableReason string `json:"disable_reason"`
}

type CustomerEnableSessionResponse struct {
	Otp string `json:"otp"`
}

type FaydaApproveRequest struct {
	RiskLevel constants.RiskLevel `json:"risk_level"`
}
type SearchCustomerByCIRequest struct {
	CifOrAccountNumber string `json:"cif_or_account_number"`
}
type Address struct {
	Zone        string `json:"zone" bson:"zone"`
	Wereda      string `json:"wereda" bson:"wereda"`
	Kebele      string `json:"kebele" bson:"kebele"`
	Region      string `json:"region" bson:"region"`
	City        string `json:"city" bson:"city"`
	SubCity     string `json:"sub_city" bson:"sub_city"`
	StreetName  string `json:"street_name" bson:"street_name"`
	HouseNumber string `json:"house_number" bson:"house_number"`
}

type CustomerDetailRespons struct {
	ID             string   `json:"id" bson:"_id"`
	CustomerCode   string   `json:"customer_code" bson:"user_code"`
	FirstName      string   `json:"first_name" bson:"first_name"`
	MiddleName     string   `json:"middle_name" bson:"middle_name"`
	LastName       string   `json:"last_name" bson:"last_name"`
	FullName       string   `json:"full_name" bson:"full_name"`
	PhoneNumber    string   `json:"phone_number" bson:"phone_number"`
	Gender         string   `json:"gender" bson:"gender"`
	AccountNumbers []string `json:"account_numbers" bson:"account_numbers"`
	BranchName     string   `json:"branch_name" bson:"branch_name"`
	DistrictName   string   `json:"district_name" bson:"district_name"`
	MothersName    string   `json:"mothers_name" bson:"mothers_name"`
	Nationality    string   `json:"nationality" bson:"nationality"`
	BirthDate      string   `json:"birth_date" bson:"birth_date"`
	Address        Address  `json:"address" bson:"address"`
	MonthlyIncome  string   `json:"monthly_income" bson:"monthly_income"`
}

type LinkedAccount struct {
	AccountNumber     string `json:"account_number" bson:"account_number"`           // linked_account
	AccountHolderName string `json:"account_holder_name" bson:"account_holder_name"` // linked_account
	AccountType       string `json:"account_type" bson:"account_type"`               // linked_account
	AccountBranchName string `json:"account_branch_name" bson:"account_branch_name"` // members
	AccountBranchCode string `json:"account_branch_code" bson:"account_branch_code"` //linked_account
	IsActive          bool   `json:"is_active" bson:"is_active"`                     /// linked_account
}

type PersonalInfo struct {
	FullName       string `json:"full_name" bson:"full_name"`             // customer_kyc.kyc_data
	Gender         string `json:"gender" bson:"gender"`                   //customer_kyc.kyc_data
	PhoneNumber    string `json:"phone_number" bson:"phone_number"`       //customer_kyc.kyc_data
	Email          string `json:"email" bson:"email"`                     // members
	CustomerNumber string `json:"customer_number" bson:"customer_number"` // members
	DateOfBirth    string `json:"date_of_birth" bson:"date_of_birth"`     // customer_kyc.kyc_data
}

type CustomerDetailResponse struct {
	ID            string          `json:"id" bson:"_id"`
	LinkedAccount []LinkedAccount `json:"linked_account" bson:"linked_account"`
	PersonalInfo  PersonalInfo    `json:"personal_info" bson:"personal_info"`
}

type CustomerListResponse struct {
	ID             string `json:"id" bson:"_id"`
	UserCode       string `json:"user_code" bson:"user_code"`
	UserID         string `json:"user_id" bson:"user_id"`
	FullName       string `json:"full_name" bson:"full_name"`
	PhoneNumber    string `json:"phone_number" bson:"phone_number"`
	CustomerNumber string `json:"customer_number" bson:"customer_number"`
	BranchCode     string `json:"branch_code" bson:"branch_code"`
	Gender         string `json:"gender" bson:"gender"`
	AccountNumber  string `json:"account_number" bson:"account_number"`
	Avatar         string `json:"avatar" bson:"avatar"`
	CreatedAt      string `json:"created_at" bson:"created_at"`
	IsBlocked      bool   `json:"is_blocked" bson:"is_blocked"`
}
