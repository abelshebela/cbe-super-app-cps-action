package model

import (
	"time"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type EventDocument struct {
	ID                  bson.ObjectID              `bson:"_id,omitempty"`
	EventCode           string                     `bson:"event_code"`
	EventName           string                     `bson:"event_name"`
	EventCity           string                     `bson:"event_city"`
	AccountNumber       string                     `bson:"account_number"`
	EventVenue          string                     `bson:"event_venue"`
	RefundPolicy        []string                   `bson:"refund_policy"`
	MICSInfo            []string                   `bson:"mics_info"`
	Restriction         domain.Restriction         `bson:"restriction"`
	Ticket              []domain.Ticket            `bson:"ticket"`
	Status              domain.EventStatus         `bson:"status"`
	EventInformation    domain.EventInformation    `bson:"event_information"`
	TicketStatistics    domain.TicketStatistics    `bson:"ticket_statistics"`
	TicketInformation   domain.TicketInformation   `bson:"ticket_information"`
	MerchantInformation domain.MerchantInformation `bson:"merchant_information"`
	Enabled             bool                       `bson:"enabled"`
	IsDeleted           bool                       `bson:"is_deleted"`
	HasRestriction      bool                       `bson:"has_restriction"`
	CreatedAt           time.Time                  `bson:"created_at"`
	DeletedAt           time.Time                  `bson:"deleted_at"`
	LastModifiedAt      time.Time                  `bson:"last_modified_at"`
}
