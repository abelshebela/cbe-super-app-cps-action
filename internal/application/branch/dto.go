package branch

import (

)

type PhoneNumber struct {
	Code   string `json:"country_code"`
	Number string `json:"number"`
}

type PinLoginRequest struct {
	DeviceUUID string `json:"device_uuid" binding:"required"`
	PinCode    string `json:"pin_code" binding:"required"`
}

type SetLoginPINRequest struct {
	PhoneNumber PhoneNumber `json:"phone_number" binding:"required"`
	PinCode     string      `json:"pin_code" binding:"required"`
}
type FetchLinkedAccountRequest struct {
	AccountNumber string `json:"account_number"`
	PhoneNumber   string `json:"phone_number"`
}


type VerifyOTPAndLinkAccountRequest struct {
    OTPCode        string `json:"otpCode"`
    CustomerNumber string `json:"customer_number"`
    AccountNo      string `json:"accountNo"`
    AccountType    string `json:"accountType"`
}
// ...existing code...

type VerifyOTPAndCreateBPSRequest struct {
    OTPCode   string `json:"otpCode"`
    AccountNo string `json:"accountNo"`
}

// ...existing code...

type LinkAccountResponse struct {
    Message string `json:"message"`
}

// ...existing code...
type LookupAccountResponse struct {
	AccountBranchType  string `json:"account_branchtype"`
	AccountBranchCode  string `json:"account_branchcode"`
	AccountNumber      string `json:"account_number"`
	CustomerNumber     string `json:"customer_number"`
	CustomerName       string `json:"customer_name"`
	AccountDescription string `json:"account_description"`
	PhoneNumber        string `json:"phone_number"`
	CustomerAddress    string `json:"customer_address"`
	DebitAllowed       bool   `json:"debit_allowed"`
	CreditAllowed      bool   `json:"credit_allowed"`
	AccountType        string `json:"account_type"`
	AccountFrozen      bool   `json:"account_frozen"`
	AccountDormant     bool   `json:"account_dormant"`
	ActiveAccount      bool   `json:"active_account"`
	AccountCurrency    string `json:"account_currency"`
}

type ApproveOrRejectBPSRequest struct {
	BPSID  string  `json:"bps_id"`
	Action string  `json:"action"`
	Reason *string `json:"reason,omitempty"`
}

// type OTPResponse struct {
// 	ID          bson.ObjectID      `json:"id" bson:"_id"`
// 	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
// 	OTPCode     string             `json:"otp_code" bson:"otp_code"`
// 	DeviceUUID  string             `json:"device_uuid" bson:"device_uuid"`
// 	OTPFor      entities.OTPFor    `json:"otp_for" bson:"otp_for"`
// 	ExpiresAt   time.Time          `json:"expires_at" bson:"expires_at"`
// 	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
// 	Status      entities.OTPStatus `json:"status" bson:"status"`
// }

type OTPRequest struct {
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	Purpose     string `json:"purpose" bson:"purpose"`
	DeviceUUID  string `json:"device_uuid" bson:"device_uuid"`
}

type OTPConfirmRequest struct {
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	OTP         string `json:"otp" bson:"otp"`
	Purpose     string `json:"purpose" bson:"purpose"`
}

type Filter struct {
	// page specifies the page number
	Page int `json:"page"`
	// per page specifies the number of results per page
	PerPage int `json:"per_page"`
	// search specifies the search term
	Search string `json:"search"`

	Filters string `json:"filters"`
}

