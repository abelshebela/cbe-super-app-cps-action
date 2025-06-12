package ticket

import (
	"time"
)

type TicketStatus string

const (
	TicketBooked   TicketStatus = "BOOKED"
	TicketExpired  TicketStatus = "EXPIRED"
	TicketPaid     TicketStatus = "PAID"
	TicketRedeemed TicketStatus = "REDEEMED"
)

type RedeemerInformation struct {
	UserID        string
	FullName      string
	AccountNumber string
	PhoneNumber   string
}

type GuestInformation struct {
	FullName    string
	Gender      string
	PhoneNumber string
}

type TicketInformation struct {
	Price             string
	Category          string
	RedeemInformation RedeemerInformation
}

type MerchantInformation struct {
	UserID      string
	Name        string
	Code        string
	PhoneNumber string
	Email       string
}

type TransactionInformation struct {
	PaidAmount            uint64
	TransactionID         string
	FTNumber              string
	AdditionalInformation struct{}
}

type Ticket struct {
	ID                     string
	EventID                string
	UserID                 string
	EventCode              string
	TicketCode             string
	TicketNumber           string
	Address                string
	RedeemerInformation    RedeemerInformation
	GuestInformation       GuestInformation
	TransactionInformation TransactionInformation
	Status                 TicketStatus
	IsGuest                bool
	IsUsed                 bool
	IsPaid                 bool
	UsedAt                 time.Time
	PaidAt                 time.Time
	CreatedAt              time.Time
	RedeemedDate           time.Time
	LastModifiedAt         time.Time
}
