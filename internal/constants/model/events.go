package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Event struct {
	ID                  bson.ObjectID             `bson:"_id,omitempty"`
	EventCode           string                    `bson:"event_code"`
	EventName           string                    `bson:"event_name"`
	EventCity           string                    `bson:"event_city"`
	AccountNumber       string                    `bson:"account_number"`
	EventVenue          string                    `bson:"event_venue"`
	RefundPolicy        []string                  `bson:"refund_policy"`
	MICSInfo            []string                  `bson:"mics_info"`
	Restriction         types.Restriction         `bson:"restriction"`
	Ticket              []Ticket                  `bson:"ticket"`
	Status              constants.EventStatus     `bson:"status"`
	EventInformation    types.EventInformation    `bson:"event_information"`
	TicketStatistics    types.TicketStatistics    `bson:"ticket_statistics"`
	TicketInformation   types.TicketInformation   `bson:"ticket_information"`
	MerchantInformation types.MerchantInformation `bson:"merchant_information"`
	Enabled             bool                      `bson:"enabled"`
	IsDeleted           bool                      `bson:"is_deleted"`
	HasRestriction      bool                      `bson:"has_restriction"`
	CreatedAt           time.Time                 `bson:"created_at"`
	DeletedAt           time.Time                 `bson:"deleted_at"`
	LastModifiedAt      time.Time                 `bson:"last_modified_at"`
}
type Ticket struct {
	Name           string `json:"name" bson:"name"`
	Category       string `json:"category" bson:"category"`
	Type           string `json:"type" bson:"type"`
	Price          uint64 `json:"price" bson:"price"`
	NumberOfTicker uint8  `json:"number_of_ticker" bson:"number_of_ticker"`
}
