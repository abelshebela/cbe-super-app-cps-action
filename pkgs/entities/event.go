package entities

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/entities/enums"
	td "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/entities/type_definition"
)

type Event struct {
	ID                  string                 `json:"id" bson:"id"`
	Code                string                 `json:"code" bson:"code"`
	Name                string                 `json:"name" bson:"name"`
	Address             td.Address             `json:"address" bson:"address"`
	AccountNumber       string                 `json:"account_number" bson:"account_number"`
	EventVenue          string                 `json:"event_venue" bson:"event_venue"`
	RefundPolicy        []string               `json:"refund_policy" bson:"refund_policy"`
	MICSInfo            []string               `json:"mics_info" bson:"mics_info"`
	Restriction         td.Restriction         `json:"restriction" bson:"restriction"`
	Ticket              td.Ticket              `json:"ticket" bson:"ticket"`
	Status              enums.EventStatus      `json:"status" bson:"status"`
	EventInformation    td.EventInformation    `json:"event_information" bson:"event_information"`
	TicketStatistics    td.TicketStatistics    `json:"ticket_statistics" bson:"ticket_statistics"`
	TicketInformation   td.TicketInformation   `json:"ticket_information" bson:"ticket_information"`
	MerchantInformation td.MerchantInformation `json:"merchant_information" bson:"merchant_information"`
	Enabled             bool                   `json:"enabled" bson:"enabled"`
	IsDeleted           bool                   `json:"is_deleted" bson:"is_deleted"`
	HasRestriction      bool                   `json:"has_restriction" bson:"has_restriction"`
	CreatedAt           time.Time              `json:"created_at" bson:"created_at"`
	DeletedAt           time.Time              `json:"deleted_at" bson:"deleted_at"`
	LastModifiedAt      time.Time              `json:"last_modified_at" bson:"last_modified_at"`
}
