package ad

import (
	"fmt"
	"mime/multipart"
	"time"

	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// AdvertDate represents the start and end dates for an advert
type AdvertDate struct {
	StartedAt time.Time `form:"started_at" json:"started_at"`
	ExpiredAt time.Time `form:"expired_at" json:"expired_at"`
}

// AdvertResponse represents the response format for an advert
type AdvertResponse struct {
	ID            string           `json:"id"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	BannerImage   string           `json:"banner_image"`
	AdvertFor     entity.AdvertFor `json:"advert_for"`
	Date          AdvertDate       `json:"date"`
	Enabled       bool             `json:"enabled"`
	CreatedAt     time.Time        `json:"created_at"`
	LastUpdatedAt time.Time        `json:"last_updated_at"`
}

// AdvertRequest represents the request payload for creating or updating an advert
type AdvertRequest struct {
	ID          string                `form:"id" json:"id"`
	Title       string                `form:"title" json:"title"`
	Description string                `form:"description" json:"description"`
	BannerImage *multipart.FileHeader `form:"banner_image" json:"banner_image"`
	AdvertFor   string                `form:"advert_for" json:"advert_for"`
	Date        AdvertDate            `form:"date" json:"date"`
}

// Validate validates the AdvertRequest struct
func (c AdvertRequest) Validate(isUpdate bool) error {
	err := validation.ValidateStruct(&c,
		validation.Field(&c.Title,
			validation.When(!isUpdate, validation.Required.Error("title is required")),
			validation.Length(3, 20).Error("TITLE_LENGTH_3_TO_20"),
		),
		validation.Field(&c.Description,
			validation.When(!isUpdate, validation.Required.Error("description is required")),
			validation.Length(30, 100).Error("DESCRIPTION_LENGTH_30_TO_100"),
		),
		validation.Field(&c.AdvertFor,
			validation.When(!isUpdate, validation.Required.Error("advert for is required")),
			validation.In(string(entity.Both), string(entity.IFB), string(entity.CB)).Error("invalid advert for field"),
		),
		validation.Field(&c.BannerImage,
			validation.When(!isUpdate, validation.Required.Error("MISSING_OR_INVALID_IMAGE")),
			validation.When(c.BannerImage != nil, validation.By(func(value interface{}) error {
				file, ok := value.(*multipart.FileHeader)
				if !ok {
					return fmt.Errorf("INVALID_FILE")
				}
				if file.Size > (2 << 20) {
					return fmt.Errorf("FILE_TOO_LARGE")
				}
				return nil
			})),
		),
		validation.Field(&c.Date),
	)

	if err != nil {
		return err
	}

	// Ensure at least one field is provided for update
	if isUpdate &&
		c.Title == "" &&
		c.Description == "" &&
		c.AdvertFor == "" &&
		c.BannerImage == nil &&
		c.Date.StartedAt.IsZero() &&
		c.Date.ExpiredAt.IsZero() {
		return fmt.Errorf("NO_DATA_PROVIDED_FOR_UPDATE")
	}

	return c.Date.Validate(isUpdate)
}

func (a AdvertDate) Validate(isUpdate bool) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.StartedAt,
			validation.When(!isUpdate, validation.Required.Error("START_DATE_REQUIRED")),
			validation.By(func(value interface{}) error {
				startedAt, ok := value.(time.Time)
				if !ok {
					if !isUpdate {
						return fmt.Errorf("The start time format is invalid. Please provide a valid date and time.")
					}
					return nil // skip validation on update if not provided
				}
				if isUpdate && startedAt.IsZero() {
					return nil // skip validation on update if zero value
				}
				if startedAt.Before(time.Now().Add(10 * time.Minute)) {
					return fmt.Errorf("The start time must be at least 10 minutes from now.")
				}
				return nil
			}),
		),
		validation.Field(&a.ExpiredAt,
			validation.When(!isUpdate, validation.Required.Error("EXPIRED_DATE_REQUIRED")),
			validation.By(func(value interface{}) error {
				expiredAt, ok := value.(time.Time)
				if !ok {
					if !isUpdate {
						return fmt.Errorf("The expiry time format is invalid. Please provide a valid date and time.")
					}
					return nil // skip validation on update if not provided
				}
				if isUpdate && expiredAt.IsZero() {
					return nil // skip validation on update if zero value
				}
				if !expiredAt.After(time.Now().Add(24 * time.Hour)) {
					return fmt.Errorf("The expiry time must be at least 24 hours from now")
				}
				if !expiredAt.After(a.StartedAt) && !a.StartedAt.IsZero() {
					return fmt.Errorf("The expiry time must be later than the start time")
				}
				return nil
			}),
		),
	)
}
