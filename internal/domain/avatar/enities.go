package avatar

import "time"

type Avatar struct {
	ID             string     `json:"id" bson:"_id"`
	Avatar         string     `json:"avatar" bson:"avatar"`
	Label          string     `json:"label" bson:"label"`
	Enable         bool       `json:"enable" bson:"enable"`
	IsDeleted      bool       `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time  `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type CPSAction struct {
	ID                 string    `json:"id"`
	ActionCode         string    `json:"action_code"`
	UniqueId           string    `json:"unique_id,omitempty"`
	MakerID            string    `json:"maker_id"`
	MakerName          string    `json:"maker_name"`
	MakerPhoneNumber   string    `json:"maker_phone_number"`
	CheckerID          string    `json:"checker_id,omitempty"`
	CheckerName        string    `json:"checker_name,omitempty"`
	CheckerPhoneNumber string    `json:"checker_phone_number,omitempty"`
	Department         string    `json:"department"`
	RejectionReason    string    `json:"rejection_reason,omitempty"`
	PreviousAction     any       `json:"previos_action"`
	CurrentAction      any       `json:"current_action"`
	ActionStatus       string    `json:"action_status"`
	ActionType         string    `json:"action_type"`
	RequestAction      string    `json:"request_action"`
	CreatedAt          time.Time `json:"created_at"`
	LastModifiedAt     time.Time `json:"last_modified_at"`
	MakerActionTime    time.Time `json:"maker_action_time"`
	CheckerActionTime  time.Time `json:"checker_action_time"`
}
