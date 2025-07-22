package cpsactions

import (
	"errors"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CPSAction struct {
	ID                 string                 `json:"id"`
	ActionCode         string                 `json:"action_code"`
	UniqueID           string                 `json:"unique_id,omitempty"`
	MakerID            string                 `json:"maker_id"`
	MakerName          string                 `json:"maker_name"`
	MakerPhoneNumber   string                 `json:"maker_phone_number"`
	CheckerID          string                 `json:"checker_id,omitempty"`
	CheckerName        string                 `json:"checker_name,omitempty"`
	CheckerPhoneNumber string                 `json:"checker_phone_number,omitempty"`
	Department         string                 `json:"department"`
	RejectionReason    string                 `json:"rejection_reason,omitempty"`
	PreviousAction     any                    `json:"previous_action"`
	CurrentAction      any                    `json:"current_action"`
	ActionStatus       constant.ActionStatus  `json:"action_status"`
	ActionType         constant.ActionType    `json:"action_type"`
	RequestAction      constant.RequestAction `json:"request_action"`
	CreatedAt          time.Time              `json:"created_at"`
	LastModifiedAt     time.Time              `json:"last_modified_at"`
	MakerActionTime    time.Time              `json:"maker_action_time"`
	CheckerActionTime  *time.Time              `json:"checker_action_time"`
}

type User struct {
	UserCode    string
	FullName    string
	PhoneNumber string
}

type AuthorizeCPSAction struct {
	ActionCode        string    `json:"action_code"`
	Department        string    `json:"department,omitempty"`
	RejectionReason   string    `json:"rejection_reason"`
	CheckerUser       User      `json:"checker_user"`
	CheckerActionTime time.Time `json:"checker_action_time,omitzero"`
}

func ToDomainCPSAction(m *model.CPSAction) *CPSAction {
	var checkerActionTime time.Time
	if m.CheckerActionTime != nil {
		checkerActionTime = *m.CheckerActionTime
	}

	return &CPSAction{
		ID:                 m.ID.Hex(),
		ActionCode:         m.ActionCode,
		UniqueID:           m.UniqueId,
		MakerID:            m.MakerID,
		MakerName:          m.MakerName,
		MakerPhoneNumber:   m.MakerPhoneNumber,
		CheckerID:          m.CheckerID,
		CheckerName:        m.CheckerName,
		CheckerPhoneNumber: m.CheckerPhoneNumber,
		Department:         m.Department,
		RejectionReason:    m.RejectionReason,
		PreviousAction:     m.PreviousAction,
		CurrentAction:      m.CurrentAction,
		ActionStatus:       constant.ActionStatus(m.ActionStatus),
		ActionType:         constant.ActionType(m.ActionType),
		RequestAction:      constant.RequestAction(m.RequestAction),
		CreatedAt:          m.CreatedAt,
		LastModifiedAt:     m.LastModifiedAt,
		MakerActionTime:    m.MakerActionTime,
		CheckerActionTime:  &checkerActionTime,
	}
}

func ToModelCPSAction(d *CPSAction) (*model.CPSAction, error) {
	if d == nil {
		return nil, errors.New("domain CPSAction is nil")
	}

	var objectID bson.ObjectID
	var err error

	if d.ID != "" {
		objectID, err = bson.ObjectIDFromHex(d.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid CPSAction ID format: %w", err)
		}
	} else {
		objectID = bson.NewObjectID()
	}

	var checkerActionTime *time.Time
	if !d.CheckerActionTime.IsZero() {
		checkerActionTime = d.CheckerActionTime
	}

	return &model.CPSAction{
		ID:                 objectID,
		ActionCode:         d.ActionCode,
		UniqueId:           d.UniqueID,
		MakerID:            d.MakerID,
		MakerName:          d.MakerName,
		MakerPhoneNumber:   d.MakerPhoneNumber,
		CheckerID:          d.CheckerID,
		CheckerName:        d.CheckerName,
		CheckerPhoneNumber: d.CheckerPhoneNumber,
		Department:         d.Department,
		RejectionReason:    d.RejectionReason,
		PreviousAction:     d.PreviousAction,
		CurrentAction:      d.CurrentAction,
		ActionStatus:       string(d.ActionStatus),
		ActionType:         string(d.ActionType),
		RequestAction:      string(d.RequestAction),
		CreatedAt:          d.CreatedAt,
		LastModifiedAt:     d.LastModifiedAt,
		MakerActionTime:    d.MakerActionTime,
		CheckerActionTime:  checkerActionTime,
	}, nil
}
