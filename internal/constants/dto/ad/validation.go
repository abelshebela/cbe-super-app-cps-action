package ad

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"fmt"
	"mime/multipart"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Validate validates the AdvertRequest struct
func (c AdvertRequest) Validate(isUpdate bool) error {
	err := validation.ValidateStruct(&c,
		validation.Field(&c.Title,
			validation.When(!isUpdate, validation.Required.Error("title is required")),
			validation.Length(3, 20).Error(localization.ErrorTitleLength3To20.Code),
		),
		validation.Field(&c.Description,
			validation.When(!isUpdate, validation.Required.Error("description is required")),
			validation.Length(30, 100).Error(localization.ErrorDescriptionLength30To100.Code),
		),
		validation.Field(&c.AdvertFor,
			validation.When(!isUpdate, validation.Required.Error("advert for is required")),
			validation.In(string(constants.BOTH_ADVERT_FOR), string(constants.IFB_ADVERT_FOR), string(constants.CB_ADVERT_FOR)).Error("invalid advert for field"),
		),
		validation.Field(&c.BannerImage,
			validation.When(!isUpdate, validation.Required.Error(localization.ErrorMissingOrInvalidImage.Code)),
			validation.When(c.BannerImage != nil, validation.By(func(value interface{}) error {
				file, ok := value.(*multipart.FileHeader)
				if !ok {
					return localization.ErrorMissingOrInvalidImage
				}
				if file.Size > (2 << 20) {
					return localization.ErrorFileTooLarge
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
						return fmt.Errorf("the start time format is invalid. please provide a valid date and time")
					}
					return nil // skip validation on update if not provided
				}
				if isUpdate && startedAt.IsZero() {
					return nil // skip validation on update if zero value
				}
				if startedAt.Before(time.Now().Add(10 * time.Minute)) {
					return fmt.Errorf("the start time must be at least 10 minutes from now")
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
						return fmt.Errorf("the expiry time format is invalid. please provide a valid date and time")
					}
					return nil // skip validation on update if not provided
				}
				if isUpdate && expiredAt.IsZero() {
					return nil // skip validation on update if zero value
				}
				if !expiredAt.After(time.Now().Add(24 * time.Hour)) {
					return fmt.Errorf("the expiry time must be at least 24 hours from now")
				}
				if !expiredAt.After(a.StartedAt) && !a.StartedAt.IsZero() {
					return fmt.Errorf("the expiry time must be later than the start time")
				}
				return nil
			}),
		),
	)
}
