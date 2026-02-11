package ad

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Validate validates the AdvertRequest struct
func (c AdvertRequest) Validate(isUpdate bool) error {
	err := validation.ValidateStruct(&c,
		validation.Field(&c.Title,
			validation.When(!isUpdate, validation.Required.Error("title is required")),
			validation.Length(3, 50).Error(localization.ErrorTitleLength3To20.Message),
		),
		validation.Field(&c.AdvertFor,
			validation.When(!isUpdate, validation.Required.Error("advert for is required")),
			validation.In(string(constants.BOTH_ADVERT_FOR), string(constants.IFB_ADVERT_FOR), string(constants.CB_ADVERT_FOR)).Error("invalid advert for field"),
		),
		validation.Field(&c.BannerImage,
			validation.When(c.BannerImage != nil, validation.By(func(value interface{}) error { return validateBannerImage(value) })),
		))

	if err != nil {
		return err
	}
	return nil
}

func validateBannerImage(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return localization.ErrorMissingOrInvalidImage
	}
	if !utils.IsValidImage(file) {
		return localization.ErrorMissingOrInvalidImage
	}

	if file.Size > (10 << 20) {
		return validation.NewError("logo", "file size exceeds 10MB limit")
	}

	return nil
}
