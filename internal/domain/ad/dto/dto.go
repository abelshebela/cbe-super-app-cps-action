package dto

import (
	"fmt"
	"mime/multipart"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AdvertFor string

const (
	IFB  AdvertFor = "IFB"
	CB   AdvertFor = "CB"
	Both AdvertFor = "ALL"
)

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

type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
)

type AdvertDate struct {
	StartedAt time.Time `form:"started_at" json:"started_at"`
	ExpiredAt time.Time `form:"expired_at" json:"expired_at"`
}

func (a AdvertDate) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.StartedAt,
			validation.Required.Error("started date is required"),
			validation.By(func(value interface{}) error {
				startedAt, ok := value.(time.Time)
				if !ok {
					return fmt.Errorf("invalid started time format")
				}

				if startedAt.Before(time.Now().Add(10 * time.Second)) {
					return fmt.Errorf("started_at must be at least 10 minutes from now")
				}
				return nil
			}),
		),
		validation.Field(&a.ExpiredAt,
			validation.Required.Error("expired at is required"),
			validation.By(func(value interface{}) error {
				expiredAt, ok := value.(time.Time)
				if !ok {
					return fmt.Errorf("invalid started time format")
				}

				if !expiredAt.After(time.Now().Add(time.Hour * 24)) {
					return fmt.Errorf("expired at must be at least 1 day after created")
				}
				return nil
			}),
		),
	)
}

type AdvertResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	BannerImage string     `json:"banner_image"`
	AdvertFor   AdvertFor  `json:"advert_for"`
	Date        AdvertDate `json:"advert_date"`
}

type CreateAdvertRequest struct {
	ID          string                `json:"id"`
	Title       string                `form:"title" json:"title"`
	Description string                `form:"description" json:"description"`
	BannerImage *multipart.FileHeader `form:"banner_image" json:"banner_image"`
	AdvertFor   AdvertFor             `form:"advert_for" json:"advert_for"`
	Date        AdvertDate            `form:"date" json:"date"`
}

func (c CreateAdvertRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Title,
			validation.Required.Error("title is required"), validation.Length(3, 10).Error("title length is between 3 and 10")),
		validation.Field(&c.Description, validation.Length(30, 100).Error("description length is between 30 and 100")),
		validation.Field(&c.AdvertFor, validation.In(
			Both,
			IFB,
			CB,
		).Error("advert value should be IFB, CB or BOTH")),
		validation.Field(&c.BannerImage, validation.By(func(value interface{}) error {
			file, ok := value.(*multipart.FileHeader)
			if !ok {
				return fmt.Errorf("invalid file")
			}
			if file.Size > (2 << 20) {
				return fmt.Errorf("file size should be less than 2MB")
			}
			return nil
		})),
		validation.Field(&c.Date),
	)
}

type UpdateAdvertRequest struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	AdvertFor   AdvertFor  `json:"advert_for"`
	Date        AdvertDate `json:"date"`
}

func (u UpdateAdvertRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Title,
			validation.NilOrNotEmpty,
			validation.Length(3, 20),
		),
	)
}
