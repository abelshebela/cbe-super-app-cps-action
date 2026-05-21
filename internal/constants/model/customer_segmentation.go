package model

import (
	"encoding/json"
	"strings"
	"time"
)

type CustomerSegmentChangeStatus string

const (
	CustomerSegmentStatusUpdated   CustomerSegmentChangeStatus = "UPDATED"
	CustomerSegmentStatusDeleted   CustomerSegmentChangeStatus = "DELETED"
	CustomerSegmentStatusUnchanged CustomerSegmentChangeStatus = "UNCHANGED"
)

func NormalizeCustomerSegmentChangeStatus(s CustomerSegmentChangeStatus) CustomerSegmentChangeStatus {
	return CustomerSegmentChangeStatus(strings.ToUpper(strings.TrimSpace(string(s))))
}

func (s CustomerSegmentChangeStatus) IsSet() bool {
	return NormalizeCustomerSegmentChangeStatus(s) != ""
}

type CustomerRoleInfo struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Name  string `json:"name" bson:"name"`
	Label string `json:"label,omitempty" bson:"label,omitempty"`
}

type CustGroupBlock struct {
	ID             string                      `json:"id,omitempty"`
	CustGroup      string                      `json:"name"`
	CustGroupLabel string                      `json:"label,omitempty"`
	Status         CustomerSegmentChangeStatus `json:"status,omitempty"`
}

type CustSegmentBlock struct {
	ID               string                      `json:"id,omitempty"`
	CustSegmentName  string                      `json:"name"`
	CustSegmentLabel string                      `json:"label,omitempty"`
	Status           CustomerSegmentChangeStatus `json:"status,omitempty"`
}

type CustSubSegmentBlock struct {
	ID                  string                      `json:"id,omitempty"`
	CustSubSegmentName  string                      `json:"name"`
	CustSubSegmentLabel string                      `json:"label,omitempty"`
	Status              CustomerSegmentChangeStatus `json:"status,omitempty"`
}

type CustomerEntry struct {
	Group      CustGroupBlock      `json:"group"`
	Segment    CustSegmentBlock    `json:"segment"`
	SubSegment CustSubSegmentBlock `json:"sub_segment"`
}

// CustSegment is the flat representation used by persistence.
// type CustSegment struct {
// 	CustGroupID    string `json:"id"`
// 	CustGroupName  string `json:"-" bson:"custGroupName"`
// 	CustGroupLabel string `json:"-" bson:"custGroupLabel,omitempty"`

// 	CustSegID    string `json:"id"`
// 	CustSegName  string `json:"-" bson:"custSegName"`
// 	CustSegLabel string `json:"-" bson:"custSegLabel,omitempty"`

// 	CustSubSegID    string `json:"id"`
// 	CustSubSegName  string `json:"-" bson:"custSubSegName"`
// 	CustSubSegLabel string `json:"-" bson:"custSubSegLabel,omitempty"`
// }

func (e CustomerEntry) ToCustSegment() CustomerEntry {
	return CustomerEntry{
		Group: CustGroupBlock{
			CustGroup:      e.Group.CustGroup,
			CustGroupLabel: e.Group.CustGroupLabel,
		},
		Segment: CustSegmentBlock{
			CustSegmentName:  e.Segment.CustSegmentName,
			CustSegmentLabel: e.Segment.CustSegmentLabel,
		},
		SubSegment: CustSubSegmentBlock{
			CustSubSegmentName:  e.SubSegment.CustSubSegmentName,
			CustSubSegmentLabel: e.SubSegment.CustSubSegmentLabel,
		},
	}
}

func CustSegmentToEntry(s CustomerEntry) CustomerEntry {
	return CustomerEntry{
		Group: CustGroupBlock{
			ID:             s.Group.ID,
			CustGroup:      s.Group.CustGroup,
			CustGroupLabel: s.Group.CustGroupLabel,
		},
		Segment: CustSegmentBlock{
			ID:               s.Segment.ID,
			CustSegmentName:  s.Segment.CustSegmentName,
			CustSegmentLabel: s.Segment.CustSegmentLabel,
		},
		SubSegment: CustSubSegmentBlock{
			ID:                  s.SubSegment.ID,
			CustSubSegmentName:  s.SubSegment.CustSubSegmentName,
			CustSubSegmentLabel: s.SubSegment.CustSubSegmentLabel,
		},
	}
}

func CustomerEntriesToCustSegments(entries []CustomerEntry) []CustomerEntry {
	segments := make([]CustomerEntry, 0, len(entries))
	for _, e := range entries {
		seg := e.ToCustSegment()
		segments = append(segments, seg)
	}
	return segments
}

func CustSegmentsToCustomerEntries(segments []CustomerEntry) []CustomerEntry {
	entries := make([]CustomerEntry, 0, len(segments))
	for _, s := range segments {
		entries = append(entries, CustSegmentToEntry(s))
	}
	return entries
}

type CustomerSegmentation struct {
	ID           string           `json:"id" bson:"_id,omitempty"`
	CustomerRole CustomerRoleInfo `json:"customer_role" bson:"customer_role"`
	Customer     []CustomerEntry  `json:"customer"`
	IsEnabled    bool             `json:"is_enabled" bson:"is_enabled"`
	IsDeleted    bool             `json:"is_deleted" bson:"is_deleted"`
	CreatedAt    time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at" bson:"updated_at"`
}

func (c *CustomerSegmentation) SyncSegmentsFromCustomer() {
	if len(c.Customer) == 0 {
		return
	}
	c.Customer = CustomerEntriesToCustSegments(c.Customer)
}

func (c *CustomerSegmentation) SyncCustomerFromSegments() {
	if len(c.Customer) == 0 {
		return
	}
	c.Customer = CustSegmentsToCustomerEntries(c.Customer)
}

func (c *CustomerSegmentation) UnmarshalJSON(data []byte) error {
	type customerSegmentation CustomerSegmentation
	aux := customerSegmentation{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*c = CustomerSegmentation(aux)
	if len(c.Customer) == 0 && len(c.Customer) > 0 {
		c.SyncSegmentsFromCustomer()
	}
	if len(c.Customer) == 0 && len(c.Customer) > 0 {
		c.SyncCustomerFromSegments()
	}
	return nil
}

func (c CustomerSegmentation) MarshalJSON() ([]byte, error) {
	if len(c.Customer) == 0 && len(c.Customer) > 0 {
		c.SyncCustomerFromSegments()
	}
	type customerSegmentation CustomerSegmentation
	return json.Marshal(customerSegmentation(c))
}
