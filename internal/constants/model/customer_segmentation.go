package model

import (
	"encoding/json"
	"time"
)

type CustomerRoleInfo struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Name  string `json:"name" bson:"name"`
	Label string `json:"label,omitempty" bson:"label,omitempty"`
}

type CustGroupBlock struct {
	CustGroup      string `json:"cust_group"`
	CustGroupLabel string `json:"cust_group_label,omitempty"`
}

type CustSegmentBlock struct {
	CustSegmentName  string `json:"cust_segment_name"`
	CustSegmentLabel string `json:"cust_segment_label,omitempty"`
}

type CustSubSegmentBlock struct {
	CustSubSegmentName  string `json:"cust_sub_segment_name"`
	CustSubSegmentLabel string `json:"cust_sub_segment_label,omitempty"`
}

type CustomerEntry struct {
	Group      CustGroupBlock      `json:"group"`
	Segment    CustSegmentBlock    `json:"segment"`
	SubSegment CustSubSegmentBlock `json:"sub_segment"`
}

// CustSegment is the flat representation used by persistence.
type CustSegment struct {
	Id              string `json:"-" bson:"id"`
	CustGroupName   string `json:"-" bson:"custGroupName"`
	CustGroupLabel  string `json:"-" bson:"custGroupLabel,omitempty"`
	CustSegName     string `json:"-" bson:"custSegName"`
	CustSegLabel    string `json:"-" bson:"custSegLabel,omitempty"`
	CustSubSegName  string `json:"-" bson:"custSubSegName"`
	CustSubSegLabel string `json:"-" bson:"custSubSegLabel,omitempty"`
}

func (e CustomerEntry) ToCustSegment() CustSegment {
	return CustSegment{
		CustGroupName:   e.Group.CustGroup,
		CustGroupLabel:  e.Group.CustGroupLabel,
		CustSegName:     e.Segment.CustSegmentName,
		CustSegLabel:    e.Segment.CustSegmentLabel,
		CustSubSegName:  e.SubSegment.CustSubSegmentName,
		CustSubSegLabel: e.SubSegment.CustSubSegmentLabel,
	}
}

func CustSegmentToEntry(s CustSegment) CustomerEntry {
	return CustomerEntry{
		Group: CustGroupBlock{
			CustGroup:      s.CustGroupName,
			CustGroupLabel: s.CustGroupLabel,
		},
		Segment: CustSegmentBlock{
			CustSegmentName:  s.CustSegName,
			CustSegmentLabel: s.CustSegLabel,
		},
		SubSegment: CustSubSegmentBlock{
			CustSubSegmentName:  s.CustSubSegName,
			CustSubSegmentLabel: s.CustSubSegLabel,
		},
	}
}

func CustomerEntriesToCustSegments(entries []CustomerEntry) []CustSegment {
	segments := make([]CustSegment, 0, len(entries))
	for _, e := range entries {
		seg := e.ToCustSegment()
		segments = append(segments, seg)
	}
	return segments
}

func CustSegmentsToCustomerEntries(segments []CustSegment) []CustomerEntry {
	entries := make([]CustomerEntry, 0, len(segments))
	for _, s := range segments {
		entries = append(entries, CustSegmentToEntry(s))
	}
	return entries
}

type CustomerSegmentation struct {
	ID                        string           `json:"id" bson:"_id,omitempty"`
	CustomerRole              CustomerRoleInfo `json:"customer_role" bson:"customer_role"`
	RemovedCustomerSegmentIDs []string         `json:"removed_customer_segment_ids,omitempty"`
	Customer                  []CustomerEntry  `json:"customer"`
	CustomerSegments          []CustSegment    `json:"-" bson:"t24_customer_sub_segments"`
	IsEnabled                 bool             `json:"is_enabled" bson:"is_enabled"`
	IsDeleted                 bool             `json:"is_deleted" bson:"is_deleted"`
	CreatedAt                 time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt                 time.Time        `json:"updated_at" bson:"updated_at"`
}

func (c *CustomerSegmentation) SyncSegmentsFromCustomer() {
	if len(c.Customer) == 0 {
		return
	}
	c.CustomerSegments = CustomerEntriesToCustSegments(c.Customer)
}

func (c *CustomerSegmentation) SyncCustomerFromSegments() {
	if len(c.CustomerSegments) == 0 {
		return
	}
	c.Customer = CustSegmentsToCustomerEntries(c.CustomerSegments)
}

func (c *CustomerSegmentation) UnmarshalJSON(data []byte) error {
	type customerSegmentation CustomerSegmentation
	aux := customerSegmentation{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*c = CustomerSegmentation(aux)
	if len(c.CustomerSegments) == 0 && len(c.Customer) > 0 {
		c.SyncSegmentsFromCustomer()
	}
	if len(c.Customer) == 0 && len(c.CustomerSegments) > 0 {
		c.SyncCustomerFromSegments()
	}
	return nil
}

func (c CustomerSegmentation) MarshalJSON() ([]byte, error) {
	if len(c.Customer) == 0 && len(c.CustomerSegments) > 0 {
		c.SyncCustomerFromSegments()
	}
	type customerSegmentation CustomerSegmentation
	return json.Marshal(customerSegmentation(c))
}
