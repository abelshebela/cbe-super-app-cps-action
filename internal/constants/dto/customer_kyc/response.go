package customerkyc

import "cbe-super-app-cps-action/internal/constants/model"

type CustomerKYCResponse struct {
	ID                   string                 `json:"id"`
	AccountType          model.AccountType      `json:"account_type"`
	CustomerCode         string                 `json:"customer_code"`
	CustomerName         CustomerInfoResp       `json:"customer_name"`
	Address              AddressResp            `json:"address"`
	Nationality          string                 `json:"nationality"`
	MaritalStatus        model.MaritalStatus    `json:"marital_status"`
	CustomerStatus       model.CustomerStatus   `json:"customer_status"`
	EmploymentStatus     model.EmploymentStatus `json:"employment_status"`
	Occupation           string                 `json:"occupation,omitempty"`
	AverageMonthlyIncome string                 `json:"average_monthly_income,omitempty"`
	EducationStatus      string                 `json:"education_status,omitempty"`
	SourceOfFund         string                 `json:"source_of_fund,omitempty"`
	KYCStatus            string                 `json:"kyc_status"`
	LivenessCheck        AlivenessCheckResp     `json:"liveness_check,omitempty"`
	MoneyLaunderingFree  bool                   `json:"money_laundering_free"`
	TermsAndConditions   string                 `json:"terms_and_conditions,omitempty"`
	CreatedAt            string                 `json:"created_at"`
	UpdatedAt            string                 `json:"updated_at"`
}

type CustomerInfoResp struct {
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name,omitempty"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	MotherName  string `json:"mother_name"`
}

type AddressResp struct {
	Country     string `json:"country"`
	Region      string `json:"region"`
	City        string `json:"city"`
	SubCity     string `json:"sub_city,omitempty"`
	Wereda      string `json:"wereda,omitempty"`
	Kebele      string `json:"kebele,omitempty"`
	HouseNumber string `json:"house_number,omitempty"`
}

type AlivenessCheckResp struct {
	IDCardFront        string `json:"id_card_front,omitempty"`
	IDCardBack         string `json:"id_card_back,omitempty"`
	LivenessCheckVideo string `json:"liveness_video,omitempty"`
}
