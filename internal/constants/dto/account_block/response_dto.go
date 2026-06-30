package accountblock

import (
	"time"

	imodel "cbe-super-app-cps-action/internal/constants/model"
)

// AccountBlockResponse represents any region, district, city, or branch
type AccountBlockResponse struct {
	ID                string                      `json:"id,omitempty"`
	Name              string                      `json:"name"`
	Code              string                      `json:"code"`
	Address           string                      `json:"address"`
	Slug              string                      `json:"slug"`
	ParentID          string                      `json:"parent_id,omitempty"`
	Parent            *imodel.AccountBlock        `json:"parent,omitempty"`
	Type              string                      `json:"type"` // R=Region, D=District, C=City, B=Branch
	RegionName        string                      `json:"region_name,omitempty"`
	FederalRegionName string                      `json:"federal_region_name,omitempty"`
	DistrictName      string                      `json:"district_name,omitempty"`
	DaoCode           int64                       `json:"dao_code,omitempty"`
	AccountType       string                      `json:"account_type,omitempty"`
	CityID            string                      `json:"city_id,omitempty"`
	RegionID          string                      `json:"region_id,omitempty"`
	DistrictID        string                      `json:"district_id,omitempty"`
	DisableReason     []imodel.AccountBlockReason `json:"disable_reason"`
	IsEnabled         bool                        `json:"is_enabled"`
	CreatedAt         time.Time                   `json:"created_at"`
	UpdatedAt         time.Time                   `json:"updated_at"`
}

// PreviousDisableReasonsResponse is returned by GET /account_block/previous_reasons/{id}.
type PreviousDisableReasonsResponse struct {
	AccountBlockID string                      `json:"account_block_id"`
	DisableReason  []imodel.AccountBlockReason `json:"disable_reason"`
}

type AccountBlockActionResponse struct {
	ID                  string        `json:"id,omitempty"`
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

// // BranchResponse represents a branch entity
// type BranchResponse struct {
// 	ID            string    `json:"id,omitempty"`
// 	BranchCode    string    `json:"branch_code"`
// 	BranchName    string    `json:"branch_name"`
// 	BranchAddress string    `json:"branch_address"`
// 	DistrictCode  string    `json:"district_code"`
// 	DistrictName  string    `json:"district_name"`
// 	RegionName    string    `json:"region_name"`
// 	RecordStat    string    `json:"record_stat"`
// 	CreatedAt     time.Time `json:"created_at"`
// 	UpdatedAt     time.Time `json:"updated_at"`
// 	Enabled       bool      `json:"enabled"`
// }

// // RegionResponse represents a region entity
// type RegionResponse struct {
// 	ID            string    `json:"id,omitempty"`
// 	RegionCode    string    `json:"region_code"`
// 	RegionName    string    `json:"region_name"`
// 	RegionAddress string    `json:"region_address"`
// 	CreatedAt     time.Time `json:"created_at"`
// 	UpdatedAt     time.Time `json:"updated_at"`
// 	Enabled       bool      `json:"enabled"`
// }

// // DistrictResponse represents a district entity
// type DistrictResponse struct {
// 	ID              string    `json:"id,omitempty"`
// 	DistrictCode    string    `json:"district_code"`
// 	DistrictName    string    `json:"district_name"`
// 	DistrictAddress string    `json:"district_address"`
// 	RegionID        string    `json:"region_id"`
// 	RegionName      string    `json:"region_name"`
// 	CreatedAt       time.Time `json:"created_at"`
// 	UpdatedAt       time.Time `json:"updated_at"`
// 	Enabled         bool      `json:"enabled"`
// }

// // CityResponse represents a city entity
// type CityResponse struct {
// 	ID           string    `json:"id,omitempty"`
// 	CityCode     string    `json:"city_code"`
// 	CityName     string    `json:"city_name"`
// 	CityAddress  string    `json:"city_address"`
// 	DistrictID   string    `json:"district_id"`
// 	DistrictName string    `json:"district_name"`
// 	RegionID     string    `json:"region_id"`
// 	RegionName   string    `json:"region_name"`
// 	CreatedAt    time.Time `json:"created_at"`
// 	UpdatedAt    time.Time `json:"updated_at"`
// 	Enabled      bool      `json:"enabled"`
// }
