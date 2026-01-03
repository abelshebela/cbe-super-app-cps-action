package model

// type CustomerKYC struct {
// 	ID                   bson.ObjectID       `json:"id" bson:"_id"`
// 	UserID               string              `json:"user_id" bson:"user_id"`
// 	KYCData              types.KYCData       `json:"kyc_data" bson:"kyc_data"`
// 	KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed" bson:"kyc_reject_reason_failed"`
// 	KYCStatus            constants.KYCStatus `json:"kyc_status" bson:"kyc_status"`
// 	KYCRejectReason      string              `json:"kyc_reject_reason" bson:"kyc_reject_reason"`
// 	KYCApproved          bool                `json:"kyc_approved" bson:"kyc_approved"`
// 	KYCActivityBy        any                 `json:"kyc_activity_by" bson:"kyc_activity_by"`
// 	KYCLevel             uint8               `json:"kyc_level" bson:"kyc_level"`
// 	Enabled              bool                `json:"enabled" bson:"enabled"`
// 	IsDeleted            bool                `json:"is_deleted" bson:"is_deleted"`
// 	CreatedAt            time.Time           `json:"created_at" bson:"created_at"`
// 	LastModifiedAt       time.Time           `json:"last_modified_at" bson:"last_modified_at"`
// 	DeletedAt            time.Time           `json:"deleted_at" bson:"deleted_at"`
// }

type AccountType string
type CustomerStatus string
type KYCStatus string
type MartialStatus string

const (
	CustomerActive  CustomerStatus = "ACTIVE"
	CustomerPending CustomerStatus = "PENDING"
	CustomerExpired CustomerStatus = "EXPIRED"
)

const (
	ConventionalAccount AccountType = "CONVENTIONAL"
	CBENoor             AccountType = "CBENOOR"
)

const (
	KYCApproved KYCStatus = "APPROVED"
	KYCPending  KYCStatus = "PENDING"
	KYCRejected KYCStatus = "REJECTED"
)

const ()

type CustomerInfo struct {
	FirstName   string
	MiddleName  string
	LastName    string
	PhoneNumber string
	email       string
	DateOfBirth string
	Gender      string
	MotherName  string
}

type Address struct {
	Country     string
	Region      string
	City        string
	SubCity     string
	Wereda      string
	Kebele      string
	HouseNumber string
}

type CustomerKYC struct {
	ID             string
	AccountType    AccountType
	CustomerCode   string
	CustomerName   CustomerInfo
	Address        Address
	MartialStatus  MartialStatus
	CustomerStatus CustomerStatus
	KYCStatus      KYCStatus
}
