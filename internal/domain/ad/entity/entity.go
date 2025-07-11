package entity

import (
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AdvertFor string

const (
	IFB  AdvertFor = "IFB"
	CB   AdvertFor = "CB"
	Both AdvertFor = "ALL"
)

type AdvertDate struct {
	StartedAt time.Time `json:"started_at" bson:"started_at"`
	ExpiredAt time.Time `json:"expired_at" bson:"expired_at"`
}

type Advert struct {
	ID            bson.ObjectID `json:"_id" bson:"_id"`
	Title         string        `json:"title" bson:"title"`
	Description   string        `json:"description" bson:"description"`
	BannerImage   string        `json:"banner_image" bson:"banner_image"`
	AdvertFor     AdvertFor     `json:"advert_for" bson:"advert_for"`
	Date          AdvertDate    `json:"advert_date" bson:"advert_date"`
	Enabled       bool          `json:"enabled" bson:"enabled"`
	IsDeleted     bool          `json:"is_deleted" bson:"is_deleted"`
	DeletedAt     time.Time     `json:"deleted_at" bson:"deleted_at"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
	LastUpdatedAt time.Time     `json:"last_updated_at" bson:"last_updated_at"`
}

type CreateAdvert struct {
	ID          string                `json:"id"`
	Title       string                `form:"title" json:"title"`
	Description string                `form:"description" json:"descritption"`
	BannerImage *multipart.FileHeader `form:"banner_image" json:"banner_image"`
	AdvertFor   AdvertFor             `form:"advert_for" json:"advert_for"`
	Date        AdvertDate            `form:"date" json:"date"`
}

type CPSAction struct {
	ID                 string
	ActionCode         string
	UniqueId           string
	MakerID            string
	MakerName          string
	MakerPhoneNumber   string
	CheckerID          string
	CheckerName        string
	CheckerPhoneNumber string
	Department         string
	RejectionReason    *string
	PreviosAction      interface{}
	CurrentAction      interface{}
	ActionStatus       ActionStatus
	ActionType         ActionType
	RequestAction      RequestAction
	CreatedAt          time.Time
	LastModifiedAt     time.Time
	MakerActionTime    time.Time
	CheckerActionTime  time.Time
}
type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
)

type User struct {
	UserCode    string `json:"user_code,omitempty" bson:"user_code"`
	FullName    string `json:"full_name,omitempty" bson:"full_name"`
	PhoneNumber string `json:"phone_number,omitempty" bson:"phone_number"`
}

type ActionStatus string

const (
	ActionPending  ActionStatus = "PENDING"
	ActionApproved ActionStatus = "APPROVED"
	ActionRejected ActionStatus = "REJECTED"
)

type RequestAction string

const (
	RequestCreateAdvert  RequestAction = "CREATE_ADVERT"
	RequestUpdateAdvert  RequestAction = "UPDATE_ADVERT"
	RequestEnableAdvert  RequestAction = "ENABLE_ADVERT"
	RequestDisableAdvert RequestAction = "DISABLE_ADVERT"
	RequestDeleteAdvert  RequestAction = "DELETE_ADVERT"
)

type AdvertResponse struct {
	Page   int       `json:"page"`
	Advert []*Advert `json:"advert"`
	Limit  int       `json:"limit"`
	Total  int64     `json:"total"`
}

type CPSActionResponse struct {
	ID                string        `json:"id,omitempty"`
	ActionCode        string        `json:"action_code,omitempty"`
	CheckerUser       User          `json:"checker_user"`
	MakerUser         User          `json:"maker_user"`
	RejectedReason    string        `json:"rejected_reason,omitempty"`
	Department        string        `json:"department,omitempty"`
	Status            ActionStatus  `json:"status,omitempty"`
	RequestAction     RequestAction `json:"request_action,omitempty"`
	ActionType        ActionType    `json:"action_type,omitempty"`
	ActionData        Advert        `json:"action_data"`
	PreviousData      any           `json:"previous_action,omitempty"`
	CurrentData       any           `json:"current_action,omitempty"`
	MakerActionTime   time.Time     `json:"maker_action_time,omitzero"`
	CheckerActionTime time.Time     `json:"checker_action_time,omitzero"`
}

type CPSActionRequest struct {
	ID                string        `json:"id,omitempty"`
	ActionCode        string        `json:"action_code,omitempty"`
	CheckerUser       User          `json:"checker_user"`
	MakerUser         User          `json:"maker_user"`
	RejectedReason    string        `json:"rejected_reason,omitempty"`
	Department        string        `json:"department,omitempty"`
	Status            ActionStatus  `json:"status,omitempty"`
	RequestAction     RequestAction `json:"request_action,omitempty"`
	ActionType        ActionType    `json:"action_type,omitempty"`
	ActionData        Advert        `json:"action_data"`
	PreviousData      any           `json:"previous_action,omitempty"`
	CurrentData       any           `json:"current_action,omitempty"`
	MakerActionTime   time.Time     `json:"maker_action_time,omitzero"`
	CheckerActionTime time.Time     `json:"checker_action_time,omitzero"`
}
