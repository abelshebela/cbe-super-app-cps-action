package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Event struct {
	ID                  bson.ObjectID             `bson:"_id,omitempty" json:"_id"`
	EventCode           string                    `bson:"event_code" json:"event_code"`
	EventName           string                    `bson:"event_name" json:"event_name"`
	EventCity           string                    `bson:"event_city" json:"event_city"`
	AccountNumber       string                    `bson:"account_number" json:"account_number"`
	EventVenue          string                    `bson:"event_venue" json:"event_venue"`
	RefundPolicy        []string                  `bson:"refund_policy" json:"refund_policy"`
	MICSInfo            []string                  `bson:"mics_info" json:"mics_info"`
	Restriction         types.Restriction         `bson:"restriction" json:"restriction"`
	Ticket              []types.Ticket            `bson:"ticket" json:"ticket"`
	Status              constants.EventStatus     `bson:"status" json:"status"`
	EventInformation    types.EventInformation    `bson:"event_information" json:"event_information"`
	TicketStatistics    types.TicketStatistics    `bson:"ticket_statistics" json:"ticket_statistics"`
	TicketInformation   types.TicketInformation   `bson:"ticket_information" json:"ticket_information"`
	MerchantInformation types.MerchantInformation `bson:"merchant_information" json:"merchant_information"`
	Enabled             bool                      `bson:"enabled" json:"enabled"`
	IsDeleted           bool                      `bson:"is_deleted" json:"is_deleted"`
	HasRestriction      bool                      `bson:"has_restriction" json:"has_restriction"`
	CreatedAt           time.Time                 `bson:"created_at" json:"created_at"`
	DeletedAt           time.Time                 `bson:"deleted_at" json:"deleted_at"`
	LastModifiedAt      time.Time                 `bson:"last_modified_at" json:"last_modified_at"`
}
