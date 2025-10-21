package ad

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"fmt"
	"mime/multipart"

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
	)

	if err != nil {
		return err
	}

	// Ensure at least one field is provided for update
	if isUpdate &&
		c.Title == "" &&
		c.Description == "" &&
		c.AdvertFor == "" &&
		c.BannerImage == nil {
		return fmt.Errorf("NO_DATA_PROVIDED_FOR_UPDATE")
	}

	return nil
}
