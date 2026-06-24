package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserActionResponsibility string
type ActionType string

const (
	MAKER   UserActionResponsibility = "MAKER"
	CHECKER UserActionResponsibility = "CHECKER"
	AUDITOR UserActionResponsibility = "AUDITOR"
)

const (
	CPSActions ActionType = "CPS_ACTION"
	BPSActions ActionType = "BPS_ACTION"
)

type UserActionLog struct {
	ID                         bson.ObjectID            `bson:"_id,omitempty" json:"id"`
	ActionID                   bson.ObjectID            `bson:"action_id" json:"action_id"`
	ActionCode                 string                   `bson:"action_code" json:"action_code"`
	GivenActionStatus          string                   `bson:"given_action_status" json:"given_action_status"`
	GivenAuditorStatus         AuditorMark              `bson:"given_auditor_status" json:"given_auditor_status"`
	ActionAuditorStatus        string                   `bson:"action_auditor_status" json:"action_auditor_status"`
	RequestAction              constants.RequestAction  `bson:"request_action" json:"request_action"`
	ActionTakenServiceName     string                   `bson:"action_taken_service_name" json:"action_taken_service_name"`
	ActionTakenServiceUniqueID string                   `bson:"action_taken_service_unique_id" json:"action_taken_service_unique_id"`
	CheckerLevel               string                   `bson:"checker_level" json:"checker_level"`
	AuditorLevel               string                   `bson:"auditor_level" json:"auditor_level"`
	UserID                     bson.ObjectID            `bson:"user_id" json:"user_id"`
	Username                   string                   `bson:"username" json:"username"`
	UserPhone                  string                   `bson:"user_phone" json:"user_phone"`
	UserRoleCode               string                   `bson:"user_role_code" json:"user_role_code"`
	UserActionResponsibilities UserActionResponsibility `bson:"user_action_responsibilities" json:"user_action_responsibilities"`
	ActionType                 ActionType               `bson:"action_type" json:"action_type"`
	IsDeleted                  bool                     `bson:"is_deleted" json:"is_deleted"`
	CreatedAt                  time.Time                `bson:"created_at" json:"created_at"`
	DeletedAt                  time.Time                `bson:"deleted_at" json:"deleted_at"`
	LastModifiedAt             time.Time                `bson:"last_modified_at" json:"last_modified_at"`
	AuditorCustomerBared       bool                     `bson:"auditor_customer_bared,omitempty" json:"auditor_customer_bared,omitempty"`
	IsCustomerReinstated       bool                     `bson:"is_customer_reinstated,omitempty" json:"is_customer_reinstated,omitempty"`
	ReinstateReason            string                   `bson:"reinstate_reason,omitempty" json:"reinstate_reason,omitempty"`
}

type LevelClaimPair struct {
	Level string // "1", "2", "3" …
	Claim string // "MARKEDASRIGHT" or "MARKEDASWRONG"
}

type UserActionLogActionCodeFilter struct {
	RequestActions        []string         `json:"request_actions"`
	ActionStatuses        []string         `json:"action_statuses"`
	AuditorStatuses       []string         `json:"auditor_statuses"`
	ActionAuditorStatuses []string         `json:"action_auditor_statuses"`
	PrivateUserIDs        []string         `json:"private_user_ids"`
	CustomerBarred        []string         `json:"customer_barred"` // Filter by customer barred status: "true", "false"
	Levels                []string         `json:"levels"`          // Generic levels - matches checker_level OR auditor_level
	CheckerLevels         []string         `json:"checker_levels"`  // Specific to checker_level field
	AuditorLevels         []string         `json:"auditor_levels"`  // Specific to auditor_level field
	Services              []string         `json:"services"`
	CheckerUserIDs        []string         `json:"checker_user_ids"`
	AuditorUserIDs        []string         `json:"auditor_user_ids"`
	MakerUserIDs          []string         `json:"maker_user_ids"`
	MakerUsernames        []string         `json:"maker_usernames"`   // Filter by maker username (case-insensitive regex)
	CheckerUsernames      []string         `json:"checker_usernames"` // Filter by checker username (case-insensitive regex)
	AuditorUsernames      []string         `json:"auditor_usernames"` // Filter by auditor username (case-insensitive regex)
	Responsibilities      []string         `json:"responsibilities"`
	LevelClaimPairs       []LevelClaimPair `json:"level_claim_pairs"`                // For auditor: level + given_auditor_status pairs
	CheckerLevelStatuses  []LevelClaimPair `json:"checker_level_statuses"`           // For checker: level + status pairs (e.g., level_1 + APPROVED)
	Search                string           `json:"search"`                           // General search across level, service, username fields
	AuditorCustomerBared  *bool            `json:"auditor_customer_bared,omitempty"` // Filter by customer barred status (nil = no filter, true/false = specific value)
}
