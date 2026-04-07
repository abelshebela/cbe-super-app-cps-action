package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CPSAction struct {
	ID                  bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ActionCode          string        `bson:"action_code" json:"action_code"`
	UniqueId            string        `bson:"unique_id" json:"unique_id,omitempty"`
	MakerID             string        `bson:"maker_id" json:"maker_id"`
	MakerName           string        `bson:"maker_name" json:"maker_name"`
	MakerPhoneNumber    string        `bson:"maker_phone_number" json:"maker_phone_number"`
	CheckerUsers        []Checker     `bson:"checker_users" json:"checker_users"`
	AuditorUsers        []Auditor     `bson:"auditor_users" json:"auditor_users"`
	AuditorCount        int32         `bson:"auditor_count" json:"auditor_count"`
	AuditorStatus       AuditorStatus `bson:"auditor_status" json:"auditor_status"`
	CurrentAuditorIndex float64       `bson:"current_auditor_index" json:"current_auditor_index"`
	CheckerCount        int32         `bson:"checker_count" json:"checker_count"`
	CurrentCheckerIndex float64       `bson:"current_checker_index" json:"current_checker_index"`
	RoleCode            string        `bson:"role_code" json:"role_code"`
	RejectionReason     string        `bson:"rejection_reason" json:"rejection_reason,omitempty"`
	CanceledReason      string        `bson:"canceled_reason" json:"canceled_reason,omitempty"`
	PreviousAction      interface{}   `bson:"previous_action" json:"previous_action,omitempty"`
	CurrentAction       interface{}   `bson:"current_action" json:"current_action,omitempty"`
	ActionStatus        string        `bson:"action_status" json:"action_status,omitempty"`
	ActionType          string        `bson:"action_type" json:"action_type,omitempty"`
	IsDeleted           bool          `bson:"is_deleted" json:"is_deleted,omitempty"`
	RequestAction       string        `bson:"request_action" json:"request_action"`
	Version             int64         `json:"version" bson:"version"`
	ReversedByRoleID    string        `bson:"reversed_by_role_id" json:"reversed_by_role_id,omitempty"`
	ReversedByID        string        `bson:"reversed_by_id" json:"reversed_by_id,omitempty"`
	ReversedByName      string        `bson:"reversed_by_name" json:"reversed_by_name,omitempty"`
	ReversedAt          time.Time     `bson:"reversed_at" json:"reversed_at,omitempty"`
	CreatedAt           time.Time     `bson:"created_at" json:"created_at,omitempty"`
	LastModifiedAt      time.Time     `bson:"last_modified_at" json:"last_modified_at,omitempty"`
	MakerActionTime     time.Time     `bson:"maker_action_time" json:"maker_action_time,omitempty"`
}

type AuditorStatus string
type AuditorMark string

const (
	MARKEDASRIGHT AuditorMark = "MARKEDASRIGHT"
	MARKEDASWRONG AuditorMark = "MARKEDASWRONG"
)

const (
	AUDITORNOTCHECKED AuditorStatus = "NOTCHECKED"
	AUDITORINPROGRESS AuditorStatus = "INPROGRESS"
	AUDITORCHECKED    AuditorStatus = "CHECKED"
)

type CheckCPSAction struct {
	UserCode      string
	FullName      string
	PhoneNumber   string
	Department    string
	RequestAction string
}

type Checker struct {
	CheckerID          string    `bson:"checker_id" json:"checker_id,omitempty"`
	RoleID             string    `bson:"role_id" json:"role_id,omitempty"`
	CheckerIndex       int32     `bson:"checker_index" json:"checker_index,omitempty"`
	CheckerName        string    `bson:"checker_name" json:"checker_name,omitempty"`
	CheckerPhoneNumber string    `bson:"checker_phone_number" json:"checker_phone_number,omitempty"`
	ApprovedAt         time.Time `bson:"approved_at" json:"approved_at,omitempty"`
}

type Auditor struct {
	AuditorID          string      `bson:"auditor_id" json:"auditor_id,omitempty"`
	RoleID             string      `bson:"role_id" json:"role_id,omitempty"`
	AuditorIndex       int32       `bson:"auditor_index" json:"auditor_index,omitempty"`
	AuditorName        string      `bson:"auditor_name" json:"auditor_name,omitempty"`
	AuditorPhoneNumber string      `bson:"auditor_phone_number" json:"auditor_phone_number,omitempty"`
	AuditorReason      string      `bson:"auditor_reason" json:"auditor_reason"`
	AuditorMark        AuditorMark `bson:"auditor_mark" json:"auditor_mark,omitempty"`
	ApprovedAt         time.Time   `bson:"approved_at" json:"approved_at,omitempty"`
}
