package constants

import "time"

type ContextKey string
type Platform string

const (
	// MESSAGE GROUP
	UpdateApp           = "Update your app"
	DeviceFound         = "Device Successfuly Found"
	OTPMessage          = "Your device lookup OTP is: %s. Valid for %d minutes."
	PhoneFound          = "Phone successfuly found"
	ImageUploadSuccess  = "image uploaded successfully"
	ProfileUploadSucess = "Profile theme set successfuly"
	PINResetSuccess     = "OTP verified for PIN reset"

	// constant
	Android                  Platform = "ANDROID"
	Ios                      Platform = "IOS"
	Prelogin                          = "PRE_LOGIN"
	OTPLength                         = 6
	DEV                               = "dev"
	UAT                               = "uat"
	Password                          = "PASSWORD"
	Login                             = "LOGIN"
	Change                            = "CHANGE"
	Permanent                         = "PERMANENT"
	TokenType                         = "TOKEN_TYPE"
	Token                             = "TOKEN"
	OTP                               = "OTP"
	VerifyOtp                         = "VERIFY_OTP"
	DeviceLookUp                      = "DEVICE_LOOKUP"
	Register                          = "REGISTER"
	Empty                             = ""
	SetPin                            = "SET_PIN"
	Pin                               = "PIN"
	OTPForRegistration                = "REGISTRATION"
	Incomplete                        = "INCOMPLETE"
	ForgetPinVerifyOtp                = "FORGET_PIN_VERIFY_OTP"
	ResetPin                          = "RESET_PIN"
	Completed                         = "COMPLETED"
	ProfileTemp                       = "PROFILE-*.TMP"
	BucketUserProfilePicture          = "USER-PROFILE-PICTURES"
	OtpExpirationTime                 = 3 * time.Minute
)

type Realm string

const (
	ELST_REALM     Realm = "ELST"
	BANK_REALM     Realm = "BANK"
	DISTRICT_REALM Realm = "DISTRICT"
	BRANCH_REALM   Realm = "BRANCH"
	MERCHANT_REALM Realm = "MERCHANT"
	COMPANY_REALM  Realm = "COMPANY"
	MEMBER_REALM   Realm = "MEMBER"
)

type Gender string

const (
	Male   Gender = "MALE"
	Female Gender = "FEMALE"
)

type KYCStatus string

const (
	KYCStatusPending  KYCStatus = "PENDING"
	KYCStatusApproved KYCStatus = "APPROVED"
	KYCStatusRejected KYCStatus = "REJECTED"
)

type BPSStatus string

const (
	BPSStatusAuthorized BPSStatus = "AUTHORIZED"
	BPSStatusDenied     BPSStatus = "DENIED"
	BPSStatusPending    BPSStatus = "PENDING"
	BPSStatusInitiated  BPSStatus = "INITIATED"
)

type MaritalStatus string

const (
	Single   MaritalStatus = "SINGLE"
	Married  MaritalStatus = "MARRIED"
	Divorced MaritalStatus = "DIVORCED"
	Widow    MaritalStatus = "WIDOW"
)

type DeviceStatus string

const (
	Linked   DeviceStatus = "LINKED"
	UnLinked DeviceStatus = "UNLINKED"
)

type AccountStatus string

const (
	Active   AccountStatus = "ACTIVE"
	InActive AccountStatus = "INACTIVE"
)

type MemberType string

const (
	CBT  MemberType = "CB"
	IFBT MemberType = "IFB"
)

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
	OTPForForgetPin        OTPFor = "FORGET_PIN"
)

type OTPStatus string

const (
	Pending  OTPStatus = "PENDING"
	Verified OTPStatus = "VERIFIED"
	Denied   OTPStatus = "DENIED"
)
