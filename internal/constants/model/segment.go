package model

import "time"

type Segment struct {
	ID                      string    `json:"id"`
	CustomerGroup           string    `json:"customer_group"`
	CustomerGroupLabel      string    `json:"customer_group_label"`
	CustomerSegment         string    `json:"customer_segment"`
	CustomerSegmentLabel    string    `json:"customer_segment_label"`
	CustomerSubsegment      string    `json:"customer_subsegment"`
	CustomerSubsegmentLabel string    `json:"customer_subsegment_label"`
	SuperappRole            string    `json:"superapp_role"`
	SuperappRoleLabel       string    `json:"superapp_role_label"`
	CheckSum                string    `json:"-"`
	IsEnabled               bool      `json:"is_enabled"`
	CreatedAt               time.Time `json:"created_at"`
	LastModifiedAt          time.Time `json:"last_modified_at"`
}
