package customerkyc

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"mime/multipart"
)

type CreateCustomerKYCRequest struct {
	AccountType          model.AccountType      `json:"account_type" validate:"required"`
	CustomerName         CustomerInfoRequest    `json:"customer_name" validate:"required"`
	Address              AddressRequest         `json:"address" validate:"required"`
	Nationality          string                 `json:"nationality" validate:"required"`
	MaritalStatus        model.MaritalStatus    `json:"marital_status" validate:"required"`
	EmploymentStatus     model.EmploymentStatus `json:"employment_status" validate:"required"`
	Occupation           string                 `json:"occupation,omitempty"`
	AverageMonthlyIncome string                 `json:"average_monthly_income,omitempty"`
	EducationStatus      string                 `json:"education_status,omitempty"`
	SourceOfFund         string                 `json:"source_of_fund,omitempty"`
	TermsAndConditions   string                 `json:"terms_and_conditions" validate:"required"`
	LivenessCheck        LivenessCheckRequest   `json:"liveness_video"`
	VerificationResult   VerificationResult     `json:"verification_result"`
}

type CustomerInfoRequest struct {
	FirstName   string `json:"first_name" validate:"required"`
	MiddleName  string `json:"middle_name,omitempty"`
	LastName    string `json:"last_name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	DateOfBirth string `json:"date_of_birth" validate:"required"`
	Gender      string `json:"gender" validate:"required,oneof=MALE FEMALE"`
	MotherName  string `json:"mother_name" validate:"required"`
}

type AddressRequest struct {
	Country     string `json:"country" validate:"required"`
	Region      string `json:"region" validate:"required"`
	City        string `json:"city" validate:"required"`
	SubCity     string `json:"sub_city,omitempty"`
	Wereda      string `json:"wereda,omitempty"`
	Kebele      string `json:"kebele,omitempty"`
	HouseNumber string `json:"house_number,omitempty"`
}

type VerificationResult struct {
	FaceMatchScore             float64 `json:"face_match_score"`
	LivenessResult             string  `json:"liveness_result"`
	DocumentAuthenticityResult string  `json:"document_authenticity_result"`
}

// ********************************************************************************************************************** //

type UpdateKYCStatusRequest struct {
	KYCStatus string `json:"kyc_status"`
}

type CustomerInfoUpdateRequest struct {
	FirstName   *string `json:"first_name,omitempty"`
	MiddleName  *string `json:"middle_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Email       *string `json:"email,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	MotherName  *string `json:"mother_name,omitempty"`
}

type AddressUpdateRequest struct {
	Country     *string `json:"country,omitempty"`
	Region      *string `json:"region,omitempty"`
	City        *string `json:"city,omitempty"`
	SubCity     *string `json:"sub_city,omitempty"`
	Wereda      *string `json:"wereda,omitempty"`
	Kebele      *string `json:"kebele,omitempty"`
	HouseNumber *string `json:"house_number,omitempty"`
}

type LivenessCheckRequest struct {
	IDCardFront        *multipart.FileHeader `json:"id_card_front,omitempty" swaggertype:"string"`
	IDCardBack         *multipart.FileHeader `json:"id_card_back,omitempty" swaggertype:"string"`
	LivenessCheckVideo *multipart.FileHeader `json:"liveness_video,omitempty" swaggertype:"string"`
}
