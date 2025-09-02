package types

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BakerOptions struct {
	Sequential bool
	UseMutex   bool
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

type LoginPIN struct {
	PIN              string    `json:"pin" bson:"pin"`
	PINHistory       [4]string `json:"pin_history" bson:"pin_histroy"`
	LastPINCreatedAt time.Time `json:"last_pin_created_at" bson:"last_pin_created_at"`
}

type UserContext struct {
	UserCode    string
	UserID      string
	FullName    string
	PhoneNumber string
	Department  string
	BranchCode  []string
	UserRole    string
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

// ==================================

type AdvertDate struct {
	StartedAt time.Time `json:"started_at" bson:"started_at"`
	ExpiredAt time.Time `json:"expired_at" bson:"expired_at"`
}

type Tier struct {
	ID        bson.ObjectID `bson:"id" json:"-"`
	Min       uint64        `bson:"min"`
	Max       uint64        `bson:"max"`
	FeeAmount uint64        `bson:"fee_amount"`
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
	PRD    string `bson:"prd"`
	VATPRD string `bson:"vatprd"`
	SFPRD  string `bson:"sfprd"`
	TRXN   string `bson:"trxn"`
}

type GLEntry struct {
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
	ID             string               `json:"id" bson:"id"`
	BranchType     constants.BranchType `json:"branch_type" bson:"branch_type"`
	ProductCode    string               `json:"product_code" bson:"product_code"`
	VATCode        string               `json:"vat_code" bson:"vat_code"`
	ServiceFeeCode string               `json:"service_fee_code" bson:"service_fee_code"`
}

type CredentialInformation struct {
	ID            bson.ObjectID             `bson:"_id" json:"id"`
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
	ID        bson.ObjectID `bson:"id"`
	Enabled   bool          `bson:"enabled"`
	IsDeleted bool          `bson:"is_deleted"`
}

type BranchInformation struct {
	BranchCode          string `json:"branch_code"`
	BranchName          string `json:"branch_name"`
	BranchAddress       string `json:"branch_address"`
	BranchOwner         string `json:"branch_owner"`
	BranchAccountNumber string `json:"branch_account_number"`
}

type KYCInformation struct {
	Name  string `bson:"name"`
	Email string `bson:"email"`
	Phone string `bson:"phone"`
}

type KYC struct {
	Status         string         `bson:"status"`
	Representative KYCInformation `bson:"representative"`
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
	AccessListName string `json:"accessListName" bson:"accessListName"`
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
