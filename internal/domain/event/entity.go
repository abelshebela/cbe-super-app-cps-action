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
	Latitude  float64
	Longitude float64
}

type MerchantInformation struct {
	MerchantID  string
	Name        string
	PhoneNumber string
	Email       string
}

type Address struct {
	Zone        string
	Wereda      string
	Kebele      string
	Region      string
	City        string
	SubCity     string
	StreetName  string
	HouseNumber string
	Location    Location
}

type Ticket struct {
	Name           string
	Category       string
	Type           string
	Price          uint64
	NumberOfTicker uint8
}

type TicketStatistics struct {
	Category           string
	Revenue            uint64
	NumberOfSoldTicket uint64
}

type TicketInformation struct {
	TotalNumberOfTicket          uint64
	TotalNumberOfAvailableTicket uint64
	TotalNumberOFUnsoldTicket    uint64
}

type EventInformation struct {
	DueDate     time.Time
	Description string
	Cover       string
	VideoLink   string
}

type Restriction struct {
	Type        RestrictionType
	Description string
}

type Event struct {
	ID                  string
	EventCode           string
	Name                string
	City                string
	AccountNumber       string
	EventVenue          string
	RefundPolicy        []string
	MICSInfo            []string
	Restriction         Restriction
	Ticket              Ticket
	Status              EventStatus
	EventInformation    EventInformation
	TicketStatistics    TicketStatistics
	TicketInformation   TicketInformation
	MerchantInformation MerchantInformation
	Enabled             bool
	IsDeleted           bool
	HasRestriction      bool
	CreatedAt           time.Time
	DeletedAt           time.Time
	LastModifiedAt      time.Time
}

type Maker struct {
	ID          string
	FullName    string
	PhoneNumber string
	Department  string
}
