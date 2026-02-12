package hqDto

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r UpdateBlockTimeRequest) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.BlockTime, validation.Required, validation.Min(uint32(1))),
	)
	if err != nil {
		return errors.New(localization.ErrorInvalidBlockTime.Code)
	}
	return nil
}

func (r UpdateArchiveTimeRequest) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.ArchiveTime, validation.Required, validation.Min(uint32(1))),
	)
	if err != nil {
		return errors.New(localization.ErrorInvalidArchiveTime.Code)
	}
	return nil
}

func (r UpdatePasswordExpiryRequest) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.PasswordExpiry, validation.Required, validation.Min(uint32(1))),
	)
	if err != nil {
		return errors.New(localization.ErrorInvalidPasswordExpiry.Code)
	}
	return nil
}
