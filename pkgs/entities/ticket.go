package entities

import (
	"time"

	"cbe-super-app-member-users/pkgs/entities/enums"
	td "cbe-super-app-member-users/pkgs/entities/type_definition"
)

type Ticket struct {
	ID                     string                    `json:"id" bson:"id"`
	EventID                string                    `json:"event_id" bson:"event_id"`
	UserID                 string                    `json:"user_id" bson:"user_id"`
	EventCode              string                    `json:"event_code" bson:"event_code"`
	TicketCode             string                    `json:"ticket_code" bson:"ticket_code"`
	TicketNumber           string                    `json:"ticket_number" bson:"ticket_number"`
	Address                string                    `json:"address" bson:"address"`
	RedeemerInformation    td.RedeemerInformation    `json:"redeemer_information" bson:"redeemer_information"`
	GuestInformation       td.GuestInformation       `json:"guest_information" bson:"guest_information"`
	TransactionInformation td.TransactionInformation `json:"transaction_information" bson:"transaction_information"`
	Status                 enums.TicketStatus        `json:"status" bson:"status"`
	IsGuest                bool                      `json:"is_guest" bson:"is_guest"`
	IsUsed                 bool                      `json:"is_used" bson:"is_used"`
	IsPaid                 bool                      `json:"is_paid" bson:"is_paid"`
	UsedAt                 time.Time                 `json:"used_at" bson:"used_at"`
	PaidAt                 time.Time                 `json:"paid_at" bson:"paid_at"`
	CreatedAt              time.Time                 `json:"created_at" bson:"created_at"`
	RedeemedDate           time.Time                 `json:"redeemed_date" bson:"redeemed_date"`
	LastModifiedAt         time.Time                 `json:"last_modified_at" bson:"last_modified_at"`
}
