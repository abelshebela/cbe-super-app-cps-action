package event

import "time"

type EventStatus string

const (
	EventUpcomming EventStatus = "UPCOMMING"
	EventLive      EventStatus = "LIVE"
	EventClosed    EventStatus = "CLOSED"
)

type RestrictionType string

const (
	AgeRestriction RestrictionType = "AGE_RESTRICTION"
)

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
	Type        RestrictionType `json:"type" bson:"type"`
	Description string          `json:"description" bson:"description"`
}

type Event struct {
	ID                  string              `json:"id,omitempty" bson:"_id,omitempty"`
	EventCode           string              `json:"event_code" bson:"event_code"`
	EventName           string              `json:"event_name" bson:"event_name"`
	EventCity           string              `json:"event_city" bson:"event_city"`
	AccountNumber       string              `json:"account_number" bson:"account_number"`
	EventVenue          string              `json:"event_venue" bson:"event_venue"`
	RefundPolicy        []string            `json:"refund_policy" bson:"refund_policy"`
	MICSInfo            []string            `json:"mics_info" bson:"mics_info"`
	Restriction         Restriction         `json:"restriction" bson:"restriction"`
	Ticket              []Ticket            `json:"ticket" bson:"ticket"`
	Status              EventStatus         `json:"status" bson:"status"`
	EventInformation    EventInformation    `json:"event_information" bson:"event_information"`
	TicketStatistics    TicketStatistics    `json:"ticket_statistics" bson:"ticket_statistics"`
	TicketInformation   TicketInformation   `json:"ticket_information" bson:"ticket_information"`
	MerchantInformation MerchantInformation `json:"merchant_information" bson:"merchant_information"`
	Enabled             bool                `json:"enabled" bson:"enabled"`
	IsDeleted           bool                `json:"is_deleted" bson:"is_deleted"`
	HasRestriction      bool                `json:"has_restriction" bson:"has_restriction"`
	CreatedAt           time.Time           `json:"created_at" bson:"created_at"`
	DeletedAt           time.Time           `json:"deleted_at" bson:"deleted_at"`
	LastModifiedAt      time.Time           `json:"last_modified_at" bson:"last_modified_at"`
}
