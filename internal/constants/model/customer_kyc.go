package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CustomerInfo struct {
	FirstName   string `json:"first_name" bson:"first_name"`
	MiddleName  string `json:"middle_name,omitempty" bson:"middle_name,omitempty"`
	LastName    string `json:"last_name" bson:"last_name"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	Email       string `json:"email" bson:"email"`
	DateOfBirth string `json:"date_of_birth" bson:"date_of_birth"`
	Gender      string `json:"gender" bson:"gender"`
	MotherName  string `json:"mother_name" bson:"mother_name"`
}

type Address struct {
	Country     string `json:"country" bson:"country"`
	Region      string `json:"region" bson:"region"`
	City        string `json:"city" bson:"city"`
	SubCity     string `json:"sub_city" bson:"sub_city"`
	Wereda      string `json:"wereda" bson:"wereda"`
	Kebele      string `json:"kebele" bson:"kebele"`
	HouseNumber string `json:"house_number" bson:"house_number"`
}

type LivenessCheck struct {
	IDCardFront        string `json:"id_card_front" bson:"id_card_front"`
	IDCardBack         string `json:"id_card_back" bson:"id_card_back"`
	LivenessCheckVideo string `json:"liveness_video" bson:"liveness_video"`
}

type VerificationResult struct {
	FaceMatchScore             float64 `json:"face_match_score" bson:"face_match_score"`
	LivenessResult             string  `json:"liveness_result" bson:"liveness_result"`
	DocumentAuthenticityResult string  `json:"document_authenticity_result" bson:"document_authenticity_result"`
}

type CustomerKYC struct {
	ID                   bson.ObjectID              `json:"id" bson:"_id,omitempty"`
	CustomerCode         string                     `json:"customer_code" bson:"customer_code"`
	AccountType          constants.AccountType      `json:"account_type" bson:"account_type"`
	CustomerName         CustomerInfo               `json:"customer_name" bson:"customer_name"`
	Address              Address                    `json:"address" bson:"address"`
	Nationality          string                     `json:"nationality" bson:"nationality"`
	MaritalStatus        constants.MaritalStatus    `json:"marital_status" bson:"marital_status"`
	CustomerStatus       constants.CustomerStatus   `json:"customer_status" bson:"customer_status"`
	EmploymentStatus     constants.EmploymentStatus `json:"employment_status" bson:"employment_status"`
	Occupation           string                     `json:"occupation" bson:"occupation"`
	AverageMonthlyIncome string                     `json:"average_monthly_income" bson:"average_monthly_income"`
	EducationStatus      string                     `json:"education_status" bson:"education_status"`
	SourceOfFund         string                     `json:"source_of_fund" bson:"source_of_fund"`
	KYCStatus            string                     `json:"kyc_status" bson:"kyc_status"`
	LivenessCheck        LivenessCheck              `json:"liveness_check" bson:"liveness_check"`
	VerificationResult   VerificationResult         `json:"verification_result" bson:"verification_result"`
	MoneyLaunderingFree  bool                       `json:"money_laundering_free" bson:"money_laundering_free"`
	TermsAndConditions   string                     `json:"terms_and_conditions" bson:"terms_and_conditions"`
	CreatedAt            time.Time                  `json:"created_at" bson:"created_at"`
	UpdatedAt            time.Time                  `json:"updated_at" bson:"updated_at"`
}
