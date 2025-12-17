package types

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// type KYCInformation struct {
// 	Name  string `json:"name" bson:"name"`
// 	Email string `json:"email" bson:"email"`
// 	Phone string `json:"phone" bson:"phone"`
// }

// type KYC struct {
// 	Status         KYCStatus      `json:"status" bson:"status"`
// 	Representative KYCInformation `json:"representative" bson:"representative"`
// }

// type BranchInformation struct {
// 	BranchCode          string `json:"branch_code"`
// 	BranchName          string `json:"branch_name"`
// 	BranchAddress       string `json:"branch_address"`
// 	BranchOwner         string `json:"branch_owner"`
// 	BranchAccountNumber string `json:"branch_account_number"`
// }

// type MiniApps struct {
// 	ID        string `json:"id" bson:"id"`
// 	Enabled   bool   `json:"enabled" bson:"enabled"`
// 	IsDeleted bool   `json:"is_deleted" bson:"is_deleted"`
// }

type CheckMiniAppMerchant struct {
	BankAccountNumber string `json:"bank_account_number"`
	Email             string `json:"email"`
	PhoneNumber       string `json:"phone_number"`
}
type MiniAppMerchantExistOptions struct {
	ExcludeID string
}

type BakerOptions struct {
	Sequential bool
	UseMutex   bool
}

// Merchant Data

type Company struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	MerchantID     string  `json:"merchant_id"`
	BusinessType   *string `json:"business_type"` // null → pointer
	Email          *string `json:"email"`         // null → pointer
	Phone          *string `json:"phone"`         // null → pointer
	APIKey         string  `json:"api_key"`
	ParentID       int     `json:"parent_id"`
	ParentMerchant string  `json:"parent_merchant"`
}

type Branch struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	BranchID          string  `json:"branch_id"`
	BusinessType      *string `json:"business_type"` // null → pointer
	AccountNumber     string  `json:"account_number"`
	AccountHolderName string  `json:"account_holder_name"`
	APIKey            string  `json:"api_key"`
	Email             *string `json:"email"` // null
	Phone             *string `json:"phone"` // null
}

type UserAccount struct {
	ID                 int      `json:"id"`
	Name               string   `json:"name"`
	Email              string   `json:"email"`
	CompanyID          int      `json:"company_id"`
	CompanyIDs         []int    `json:"company_ids"`
	DefaultMerchantID  string   `json:"default_merchant_id"`
	AllowedMerchantIDs []string `json:"allowed_merchant_ids"`
}

// for merchant lookup end

type AccountLookupData struct {
	AccountNumber  string
	CustomerName   string
	Restriction    string
	Currency       string
	AccountType    string
	AccountStatus  string
	AccountHolder  string
	WorkingBalance string
	CustomerID     string
}

type Account struct {
	ID                 string    `json:"id"`
	AccountBranchType  string    `json:"account_branchtype"`
	AccountBranchCode  string    `json:"account_branchcode"`
	AccountNumber      string    `json:"account_number"`
	CustomerNumber     string    `json:"customer_number"`
	CustomerName       string    `json:"customer_name"`
	CustomerMotherName string    `json:"customer_mother_name"`
	PhoneNumber        string    `json:"phone_number"`
	CustomerAddress    string    `json:"customer_address"`
	AccountType        string    `json:"account_type"`
	Gender             string    `json:"gender"`
	Picture            string    `json:"picture"`
	DebitAllowed       bool      `json:"debit_allowed"`
	CreditAllowed      bool      `json:"credit_allowed"`
	AccountFrozen      bool      `json:"account_frozen"`
	AccountDormant     bool      `json:"account_dormant"`
	ActiveAccount      bool      `json:"active_account"`
	AccountCurrency    string    `json:"account_currency"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type MakerChecker struct {
	Linkers   Linkers `json:"linkers" bson:"linkers"`
	Unlinkers Linkers `json:"unlinkers" bson:"unlinkers,omitempty"`
}

type Linkers struct {
	Maker   string `json:"maker" bson:"maker"`
	Checker string `json:"checker" bson:"checker"`
}
type AccountInfo struct {
	ID                 string `json:"id"`
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

type LoginPIN struct {
	PIN              string    `json:"pin" bson:"pin"`
	PINHistory       [4]string `json:"pin_history" bson:"pin_histroy"`
	LastPINCreatedAt time.Time `json:"last_pin_created_at" bson:"last_pin_created_at"`
}

type UserContext struct {
	UserCode     string
	UserID       string
	FullName     string
	PhoneNumber  string
	Department   string
	BranchCode   []string
	UserRole     string
	CheckerIndex string
}

type RegistrationRecord struct {
	ID          string
	PhoneNumber string
	DeviceUUID  string
	Platform    string
	FullName    string
	OTP         string
	OTPFor      string
	Status      string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	Attempts    int
	MaxAttempts int
	Email       string
}

type UserInfo struct {
	UserID      string
	FullName    string
	PhoneNumber string
	Action      string
}

type AccountUser struct {
	ID               string
	PhoneNumber      string
	KYCLevel         uint8
	BranchCode       string
	FullName         string
	RegistrationType string
	AndOrStatus      bool
}

type LinkedAccount struct {
	UserID            string
	CustomerNumber    string
	AccountNumber     string
	AccountHolderName string
	AccountType       string
	BranchCode        string
	RegistrationType  string
	AndOrStatus       bool
	CurrencyCode      string
	IsMain            bool
}

type Password struct {
	Salt             string    `json:"salt" bson:"salt"`
	CurrentPassword  string    `json:"current_password" bson:"current_password"`
	OldPassword      [4]string `json:"old_password" bson:"old_password,omitempty"`
	PasswordChangeAt time.Time `json:"password_changed_at" bson:"password_changed_at"`
}

type Checker struct {
	CheckerID          string    `bson:"checker_id" json:"checker_id,omitempty"`
	RoleID             string    `bson:"role_id" json:"role_id,omitempty"`
	CheckerIndex       int32     `bson:"checker_index" json:"checker_index,omitempty"`
	CheckerName        string    `bson:"checker_name" json:"checker_name,omitempty"`
	CheckerPhoneNumber string    `bson:"checker_phone_number" json:"checker_phone_number,omitempty"`
	ApprovedAt         time.Time `bson:"approved_at" json:"approved_at,omitempty"`
}

type Auditor struct {
	AuditorID          string    `bson:"auditor_id" json:"auditor_id,omitempty"`
	RoleID             string    `bson:"role_id" json:"role_id,omitempty"`
	AuditorIndex       int32     `bson:"auditor_index" json:"auditor_index,omitempty"`
	AuditorName        string    `bson:"auditor_name" json:"auditor_name,omitempty"`
	AuditorPhoneNumber string    `bson:"auditor_phone_number" json:"auditor_phone_number,omitempty"`
	ApprovedAt         time.Time `bson:"approved_at" json:"approved_at,omitempty"`
}

// ==================================

type AdvertDate struct {
	StartedAt time.Time `json:"started_at" bson:"started_at"`
	ExpiredAt time.Time `json:"expired_at" bson:"expired_at"`
}

type Tier struct {
	// ID        bson.ObjectID `json:"id" bson:"id"`
	Min       uint64 `json:"min" bson:"min"`
	Max       uint64 `json:"max" bson:"max"`
	FeeAmount uint64 `json:"fee_amount" bson:"fee_amount"`
}

type Cap struct {
	KYCLevel           string `json:"kyc_level" bson:"kyc_level"`
	ISingleCap         uint64 `json:"individual_single_cap" bson:"individual_single_cap"`
	IDailyCap          uint64 `json:"individual_daily_cap" bson:"individual_daily_cap"`
	CorporateSingleCap uint64 `json:"corporate_single_cap" bson:"corporate_single_cap"`
	CorporateDailyCap  uint64 `json:"corporate_daily_cap" bson:"corporate_daily_cap"`
	MinAmount          uint64 `json:"min_amount" bson:"min_amount"`
}

type ProductCodes struct {
	PRD    string `bson:"prd" json:"prd"`
	VATPRD string `bson:"vatprd" json:"vatprd"`
	SFPRD  string `bson:"sfprd" json:"sfprd"`
	TRXN   string `bson:"trxn" json:"trxn"`
}

type GLEntry struct {
	ProductAccount    string `json:"product_account" bson:"product_account"`
	ProductBranchCode string `json:"product_branch_code" bson:"product_branch_code"`
	ServiceAccount    string `json:"service_account" bson:"service_account"`
	ServiceBranchCode string `json:"service_branch_code" bson:"service_branch_code"`
	VatAccount        string `json:"vat_account" bson:"vat_account"`
	VatBranchCode     string `json:"vat_branch_code" bson:"vat_branch_code"`
}

type IFBglEntry struct {
	ProductAccount    string `json:"product_account" bson:"product_account"`
	ProductBranchCode string `json:"product_branch_code" bson:"product_branch_code"`
	ServiceAccount    string `json:"service_account" bson:"service_account"`
	ServiceBranchCode string `json:"service_branch_code" bson:"service_branch_code"`
	VatAccount        string `json:"vat_account" bson:"vat_account"`
	VatBranchCode     string `json:"vat_branch_code" bson:"vat_branch_code"`
}

type Location struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type MerchantInformation struct {
	MerchantID          string `json:"merchant_id" bson:"merchant_id"`
	MercahntName        string `json:"merchant_name" bson:"merchant_name"`
	MerchantPhoneNumber string `json:"merchant_phone_number" bson:"merchant_phone_number"`
	MerchantEmail       string `json:"merchant_email" bson:"merchant_email"`
}

type Ticket struct {
	Name           string `json:"name" bson:"name"`
	Category       string `json:"category" bson:"category"`
	Type           string `json:"type" bson:"type"`
	Price          uint64 `json:"price" bson:"price"`
	NumberOfTicker uint8  `json:"number_of_ticker" bson:"number_of_ticker"`
}

type Address struct {
	Zone   string `json:"zone" bson:"zone"`
	Kebele string `json:"kebele" bson:"kebele"`
	Woreda string `json:"woreda" bson:"woreda"`
	Region string `json:"region" bson:"region"`
}

type KYCData struct {
	Sub           string           `json:"sub" bson:"sub"`
	FullName      string           `json:"full_name" bson:"full_name"`
	PhoneNumber   string           `json:"phone_number" bson:"phone_number"`
	Gender        string           `json:"gender" bson:"gender"`
	Picture       string           `json:"picture" bson:"picture"`
	SelfiePhoto   string           `json:"selfie_photo" bson:"selfie_photo"`
	Nationality   string           `json:"nationality" bson:"nationality"`
	BirthDate     time.Time        `json:"birth_date" bson:"birth_date"`
	DocumentFront string           `json:"document_front" bson:"document_front"`
	DocumentBack  string           `json:"document_back" bson:"document_back"`
	MonthlyIncome string           `json:"monthly_income" bson:"monthly_income"`
	AccountType   string           `json:"account_type" bson:"account_type"`
	Country       string           `json:"country" bson:"country"`
	MothersName   string           `json:"mothers_name" bson:"mothers_name"`
	Vendor        constants.Vendor `json:"vendor" bson:"vendor"`
	Address       Address          `json:"address" bson:"address"`
}

type TicketStatistics struct {
	Category           string `json:"category" bson:"category"`
	Revenue            uint64 `json:"revenue" bson:"revenue"`
	NumberOfSoldTicket uint64 `json:"number_of_sold_ticket" bson:"number_of_sold_ticket"`
}

type TicketInformation struct {
	TotalNumberOfTicket          uint64 `json:"total_number_of_ticket" bson:"total_number_of_ticket"`
	TotalNumberOfAvailableTicket uint64 `json:"total_number_of_available_ticket" bson:"total_number_of_available_ticket"`
	TotalNumberOFUnsoldTicket    uint64 `json:"total_number_of_unsold_ticket" bson:"total_number_of_unsold_ticket"`
}

// InAppBroadcastMessage is the payload for in-app broadcast notifications
// that will be wrapped by shared/notification/dto.NewNotificationMessage
// and sent to Kafka with type "in_app_broadcast".
type InAppBroadcastMessage struct {
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`       // should be "inapp"
	ExpiresAt time.Time `json:"expires_at"` // RFC3339 when marshaled
}

type EventInformation struct {
	StartDate   time.Time `json:"start_date" bson:"start_date"`
	DueDate     time.Time `json:"due_date" bson:"due_date"`
	Description string    `json:"description" bson:"description"`
	Cover       string    `json:"cover" bson:"cover"`
	VideoLink   string    `json:"video_link" bson:"video_link"`
}

type Restriction struct {
	Type        constants.RestrictionType `json:"type" bson:"type"`
	Description string                    `json:"description" bson:"description"`
}

type ProductCode struct {
	ID             string               `json:"_id" bson:"_id"`
	BranchType     constants.BranchType `json:"branch_type" bson:"branch_type"`
	ProductCode    string               `json:"product_code" bson:"product_code"`
	VATCode        string               `json:"vat_code" bson:"vat_code"`
	ServiceFeeCode string               `json:"service_fee_code" bson:"service_fee_code"`
}

type CredentialInformation struct {
	Environment   constants.EnvironmentType `bson:"environment" json:"environment"`
	MerchantAppID string                    `bson:"merchant_app_id" json:"merchant_app_id"`
	FabricAppID   string                    `bson:"fabric_app_id" json:"fabric_app_id"`
	ShortCode     string                    `bson:"short_code" json:"short_code"`
	AppSecret     string                    `bson:"app_secret" json:"app_secret"`
	PrivateKey    string                    `bson:"private_key" json:"private_key"`
	PublicKey     string                    `bson:"public_key" json:"public_key"`
	Timestamp     time.Time                 `bson:"timestamp" json:"timestamp"`
	Signature     string                    `bson:"signature" json:"-"`
	MiniAppCode   string                    `bson:"mini_app_code" json:"mini_app_code"`
}

type MiniApps struct {
	ID        bson.ObjectID `json:"_id" bson:"_id"`
	Enabled   bool          `json:"enabled" bson:"enabled"`
	IsDeleted bool          `json:"is_deleted" bson:"is_deleted"`
}

type BranchInformation struct {
	BranchCode          string `json:"branch_code" bson:"branch_code"`
	BranchName          string `json:"branch_name" bson:"branch_name"`
	BranchAddress       string `json:"branch_address" bson:"branch_address"`
	BranchOwner         string `json:"branch_owner" bson:"branch_owner"`
	BranchAccountNumber string `json:"branch_account_number" bson:"branch_account_number"`
}
type KYCInformation struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Phone string `json:"phone" bson:"phone"`
}

type KYC struct {
	Status         string         `json:"status" bson:"status"`
	Representative KYCInformation `json:"representative" bson:"representative"`
}

type DonationImage struct {
	ID        string    `json:"id" bson:"id"`
	PhotoURL  string    `json:"photo_url" bson:"photo_url"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type UserSearchResult struct {
	ID                 string    `json:"id"`
	AccountBranchType  string    `json:"account_branchtype"`
	AccountBranchCode  string    `json:"account_branchcode"`
	AccountNumber      string    `json:"account_number"`
	CustomerNumber     string    `json:"customer_number"`
	CustomerName       string    `json:"customer_name"`
	AccountDescription string    `json:"account_description"`
	PhoneNumber        string    `json:"phone_number"`
	CustomerAddress    string    `json:"customer_address"`
	DebitAllowed       bool      `json:"debit_allowed"`
	CreditAllowed      bool      `json:"credit_allowed"`
	AccountType        string    `json:"account_type"`
	AccountFrozen      bool      `json:"account_frozen"`
	AccountDormant     bool      `json:"account_dormant"`
	ActiveAccount      bool      `json:"active_account"`
	AccountCurrency    string    `json:"account_currency"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type User struct {
	ID          string
	FullName    string
	PhoneNumber string
	Department  string
}

type SubAccessList struct {
	Key            string `json:"key" bson:"key"`
	Enabled        bool   `json:"enabled" bson:"enabled"`
	AccessListName string `json:"access_list_ame" bson:"access_list_name"`
}

type Response struct {
	Question string      `json:"question" bson:"question"`
	Answer   interface{} `json:"answer" bson:"answer"`
	Type     string      `json:"type" bson:"type"`
	Options  interface{} `json:"options" bson:"options"`
}

type PaginationMeta struct {
	TotalDocs     int64 `json:"total_docs"`
	Limit         int   `json:"limit"`
	TotalPages    int   `json:"total_pages"`
	Page          int   `json:"page"`
	PagingCounter int   `json:"paging_counter"`
	HasPrevPage   bool  `json:"has_prev_page"`
	HasNextPage   bool  `json:"has_next_page"`
	PrevPage      *int  `json:"prev_page,omitempty"`
	NextPage      *int  `json:"next_page,omitempty"`
}

type Filter struct {
	Page    int                    `json:"page"`
	PerPage int                    `json:"per_page"`
	Search  string                 `json:"search"`
	Filters map[string]interface{} `json:"filters"`
}

type Services struct {
	Self  bool `json:"self" bson:"self"`
	Other bool `json:"other" bson:"other"`
	Agent bool `json:"agent" bson:"agent"`
}

type SMSKafkaMessage struct {
	Recipient   string `json:"recipient"`
	MessageBody string `json:"message_body"`
}

type EmailContact struct {
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	Email string `json:"email" bson:"email"`
}

type EmailKafkaMessage struct {
	Recipients         []EmailContact         `json:"recipients"`
	CC                 []EmailContact         `json:"cc,omitempty"`
	Subject            string                 `json:"subject"`
	Type               string                 `json:"type"` // "otp", "message", "transaction", "request"
	OTPCode            string                 `json:"otp_code,omitempty"`
	Receiver           string                 `json:"receiver,omitempty"`
	MessageBody        string                 `json:"message_body,omitempty"`
	Link               string                 `json:"link,omitempty"`
	CustomerName       string                 `json:"customer_name,omitempty"`
	TransactionDetails map[string]interface{} `json:"transaction_details,omitempty"`
	Priority           int                    `json:"priority,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}
