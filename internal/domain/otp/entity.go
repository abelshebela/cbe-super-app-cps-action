package otp

import "time"

type OTPFor string

const (
	OTPForLogin            OTPFor = "LOGIN"
	OTPForAddAccount       OTPFor = "ADD_ACCOUNT"
	OTPForPINSet           OTPFor = "PIN_SET"
	OTPForTransfer         OTPFor = "TRANSFER"
	OTPForAcctivateAccount OTPFor = "ACCTIVATE_ACCOUNT"
	OTPForPINReset         OTPFor = "PIN_RESET"
	OTPForSignup           OTPFor = "SIGNUP"
	OTPForAccountLink      OTPFor = "ACCOUNT_LINK"
	OTPForChangePhone      OTPFor = "CHANGE_PHONE"
	OTPForDetachPhone      OTPFor = "DETACH_PHONE"
	OTPForAttachPhone      OTPFor = "ATTACH_PHONE"
	OTPForEnable           OTPFor = "ENABLE"
	OTPForTransferLimit    OTPFor = "TRANSFER_LIMIT"
	OTPForChangeEmail      OTPFor = "CHANGE_EMAIL"
	OTPForUpgradeLimit     OTPFor = "UPGRADE_LIMIT"
)

type OTPStatus string

const (
	Pending  OTPStatus = "PENDING"
	Verified OTPStatus = "VERIFIED"
	Denied   OTPStatus = "DENIED"
)

type Realm string

const (
	ElstRealm     Realm = "ELST"
	BankRealm     Realm = "BANK"
	DistrictRealm Realm = "DISTRICT"
	BranchRealm   Realm = "BRANCH"
	MerchantRealm Realm = "MERCHANT"
	CompanyRealm  Realm = "COMPANY"
	MemberRealm   Realm = "MEMBER"
)

type OTP struct {
	ID            string    `json:"id,omitempty" bson:"_id,omitempty"`
	PhoneNumber   string    `json:"phone_number" bson:"phone_number"`
	AccountNumber *string   `json:"account_number" bson:"account_number"`
	Realm         Realm     `json:"user_realm" bson:"user_realm"`
	Email         string    `json:"email" bson:"email"`
	UserCode      string    `json:"user_code" bson:"user_code"`
	OTPCode       string    `json:"otp_code" bson:"otp_code"`
	BillNo        *string   `json:"bill_no,omitempty" bson:"bill_no,omitempty"`
	DeviceUUID    *string   `json:"device_uuid,omitempty" bson:"device_uuid,omitempty"`
	OTPFor        OTPFor    `json:"otp_for" bson:"otp_for"`
	Status        OTPStatus `json:"status" bson:"status"`
	IsDeleted     bool      `json:"is_deleted" bson:"is_deleted"`
	ExpiresAt     time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	LastModified  time.Time `json:"last_modified" bson:"last_modified"`
	DeletedAt     time.Time `json:"deleted_at" bson:"deleted_at,omitempty"`
}
