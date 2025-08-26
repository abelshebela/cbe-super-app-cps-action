package type_definition

import (
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/entities/enums"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductCode struct {
	ID          string           `json:"id" bson:"id"`
	BranchType  enums.BranchType `json:"branch_type" bson:"branch_type"`
	ProductCode string           `json:"product_code" bson:"product_code"`
}

type HealthCheck struct {
	Name        string                 `json:"name"`
	Status      enums.HealthStatus     `json:"status"`
	Message     string                 `json:"message,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	LastChecked time.Time              `json:"last_checked"`
	Duration    time.Duration          `json:"duration"`
}

type AppError struct {
	Type      enums.ErrorType `json:"type"`
	Code      string          `json:"code"`
	Message   string          `json:"message"`
	Details   string          `json:"details,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
	RequestID string          `json:"request_id,omitempty"`
	Err       error           `json:"-"`
}

type KYCInformation struct {
	Name        string `json:"name" bson:"name"`
	Email       string `json:"email" bson:"email"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}

type User struct {
	UserID      string    `json:"user_id"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	Timestamp   time.Time `json:"timestamp"`
}

type MakerAndChecker struct {
	Maker   User `json:"maker" bson:"maker"`
	Checker User `json:"checker" bson:"checker"`
}

type BranchInformation struct {
	ID            string `json:"id" bson:"id"`
	Code          string `json:"code" bson:"code"`
	Name          string `json:"name" bson:"name"`
	Address       string `json:"address" bson:"address"`
	Owner         string `json:"owner" bson:"owner"`
	AccountNumber string `json:"account_number" bson:"account_number"`
}

type AppType struct {
	UAT        string `json:"uat" bson:"uat"`
	Production string `json:"production" bson:"production"`
	Test       string `json:"test" bson:"test"`
	Dev        string `json:"dev" bson:"dev"`
}

type CredentialInformation struct {
	ID            string                `json:"id" bson:"id"`
	Environment   enums.EnvironmentType `json:"environment" bson:"environment"`
	MerchantAppID string                `json:"merchant_app_id" bson:"merchant_app_id"`
	FabricAppID   string                `json:"fabric_app_id" bson:"fabric_app_id"`
	ShortCode     string                `json:"short_code" bson:"short_code"`
	AppSecret     string                `json:"app_secret" bson:"app_secret"`
	PrivateKey    string                `json:"private_key" bson:"private_key"`
	PublicKey     string                `json:"public_key" bson:"public_key"`
}

type Address struct {
	Zone        string   `json:"zone" bson:"zone"`
	Wereda      string   `json:"wereda" bson:"wereda"`
	Kebele      string   `json:"kebele" bson:"kebele"`
	Region      string   `json:"region" bson:"region"`
	City        string   `json:"city" bson:"city"`
	SubCity     string   `json:"sub_city" bson:"sub_city"`
	StreetName  string   `json:"street_name" bson:"street_name"`
	HouseNumber string   `json:"house_number" bson:"house_number"`
	Location    Location `json:"location" bson:"location"`
}

type ActionData struct {
	UseCode     string `json:"user_code" bson:"user_code"`
	FullName    string `json:"full_name" bson:"full_name"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}

type TicketType struct {
	TicketName      string `json:"ticket_name"`
	NumberOfTickets int    `json:"number_of_tickets"`
	TicketPrice     string `json:"ticket_price"`
}

type EventDTO struct {
	EventID          string          `json:"event_id"`
	EventCode        string          `json:"event_code"`
	CoverImage       string          `json:"cover_image"`
	EventName        string          `json:"event_name"`
	EventDescription string          `json:"event_description"`
	TotalTicketCount int             `json:"total_ticket_count"`
	EventStartDate   time.Time       `json:"event_startdate"`
	EventEndDate     time.Time       `json:"event_enddate"`
	EventStatus      string          `json:"event_status"`
	EventVenue       string          `json:"event_venue"`
	EventCity        string          `json:"event_city"`
	MerchantName     string          `json:"merchant_name"`
	MerchantPhone    string          `json:"merchant_phone"`
	TicketTypes      []TicketTypeDTO `json:"ticket_types"`
	TicketSales      []TicketSaleDTO `json:"ticket_sales"`
}

// TicketTypeDTO represents a ticket type for an event
type TicketTypeDTO struct {
	TicketName  string `json:"ticket_name"`
	TicketPrice int    `json:"ticket_price"`
	TicketCount int    `json:"ticket_count"`
}

// TicketSaleDTO represents ticket sales information
type TicketSaleDTO struct {
	TicketName    string `json:"ticket_name"`
	TicketRevenue int    `json:"ticket_revenue"`
	TicketSold    int    `json:"ticket_sold"`
	TicketCount   int    `json:"ticket_count"`
}
type SurveyResponse struct {
	Question string `bson:"question" json:"question" validate:"required"`
	Answer   string `bson:"answer" json:"answer" validate:"required"`
}

type AdvertDate struct {
	StartedAt time.Time `json:"started_at" bson:"started_at"`
	ExpiredAt time.Time `json:"expired_at" bson:"expired_at"`
}

type CurrentAction struct {
	Id     []string `json:"id" bson:"id"`
	Action bool     `json:"action" bson:"action"`
}

type LoginPIN struct {
	PIN              string    `json:"pin" bson:"pin"`
	PINHistory       [4]string `json:"pin_history" bson:"pin_histroy"`
	LastPINCreatedAt time.Time `json:"last_pin_creared_at" bson:"last_pin_created_at"`
}

type DeviceLinkHistroy struct {
	ID         bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     bson.ObjectID `json:"user_id" bson:"user_id"`
	Device     Device        `json:"device" bson:"device"`
	LinkedAt   time.Time     `json:"linked_at" bson:"linked_at"`
	UnLinkedAt time.Time     `json:"unlinked_at" bson:"unlinked_at"`
}

type Device struct {
	DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
	AppVersion string `json:"app_version" bson:"app_version"`
}

type Password struct {
	Salt             string    `json:"salt" bson:"salt"`
	CurrentPassword  string    `json:"current_password" bson:"current_password"`
	OldPassword      [4]string `json:"old_password" bson:"old_password,omitempty"`
	PasswordChangeAt time.Time `json:"password_changed_at" bson:"password_changed_at"`
}

type FullName struct {
	FirstName  string `json:"first_name" bson:"first_name"`
	MiddleName string `json:"middle_name" bson:"middle_name"`
	LastName   string `json:"last_name" bson:"last_name"`
}

type PhoneNumber struct {
	Code       string `json:"code" bson:"code"`
	Number     string `json:"number" bson:"number"`
	IsVerified *bool  `json:"is_verfied,omitempty" bson:"is_verfied,omitempty"`
}

type ProductCodes struct {
	PRD    string `json:"prd" bson:"prd"`
	VATPRD string `json:"vatprd" bson:"vatprd"`
	SFPRD  string `json:"sfprd" bson:"sfprd"`
	TRXN   string `json:"trxn" bson:"trxn"`
}

type GLEntry struct {
	ProductAccount    string `json:"product_account" bson:"product_account"`
	ProductBranchCode string `json:"product_branch_code" bson:"product_branch_code"`
	ServiceAccount    string `json:"service_account" bson:"service_account"`
	ServiceBranchCode string `json:"service_branch_code" bson:"service_branch_code"`
	VatAccount        string `json:"vat_account" bson:"vat_account"`
	VatBranchCode     string `json:"vat_branch_code" bson:"vat_branch_code"`
}
type Filter struct {
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Search  string `json:"search"`
	Filters string `json:"filters"`
}
type Tier struct {
	ID        *string `json:"id" bson:"id"`
	Min       uint64  `json:"min" bson:"min"`
	Max       uint64  `json:"max" bson:"max"`
	FeeAmount uint64  `json:"fee_amount" bson:"fee_amount"`
}

type Cap struct {
	SingleCap uint64 `json:"single_cap" bson:"single_cap"`
	DailyCap  uint64 `json:"daily_cap" bson:"daily_cap"`
	MinAmount uint64 `json:"min_amount" bson:"min_amount"`
}

type Location struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type MerchantInformation struct {
	MerchantID  string `json:"merchant_id" bson:"merchant_id"`
	Name        string `json:"name" bson:"name"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	Email       string `json:"email" bson:"email"`
}

type Ticket struct {
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
	TotalNumberOfTicket          uint64              `json:"total_number_of_ticket" bson:"total_number_of_ticket"`
	TotalNumberOfAvailableTicket uint64              `json:"total_number_of_available_ticket" bson:"total_number_of_available_ticket"`
	TotalNumberOFUnsoldTicket    uint64              `json:"total_number_of_unsold_ticket" bson:"total_number_of_unsold_ticket"`
	Price                        string              `json:"price,omitempty" bson:"price,omitempty"`
	Category                     string              `json:"category,omitempty" bson:"category,omitempty"`
	RedeemInformation            RedeemerInformation `json:"redeem_information,omitempty" bson:"redeem_information,omitempty"`
}

type EventInformation struct {
	DueDate     time.Time `json:"due_date" bson:"due_date"`
	Description string    `json:"description" bson:"description"`
	Cover       string    `json:"cover" bson:"cover"`
	VideoLink   string    `json:"video_link" bson:"video_link"`
}

type Restriction struct {
	Type        string `json:"type" bson:"type"`
	Description string `json:"description" bson:"description"`
}

type RedeemerInformation struct {
	UserID        string `json:"user_id" bson:"user_id"`
	FullName      string `json:"full_name" bson:"full_name"`
	AccountNumber string `json:"account_number" bson:"account_number"`
	PhoneNumber   string `json:"phone_number" bson:"phone_number"`
}

type GuestInformation struct {
	FullName    string `json:"full_name" bson:"full_name"`
	Gender      string `json:"gender" bson:"gender"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}

type TransactionInformation struct {
	PaidAmount            uint64   `json:"paid_amount" bson:"paid_amount"`
	TransactionID         string   `json:"transaction_id" bson:"transaction_id"`
	FTNumber              string   `json:"ft_number" bson:"ft_number"`
	AdditionalInformation struct{} `json:"additional_information" bson:"additional_information"`
}
