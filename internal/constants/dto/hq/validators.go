package hqDto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r UpdateBlockTimeRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BlockTime, validation.Required, validation.Min(uint32(1)).Error("block_time must be greater than 0")),
	)
}
func (r UpdateArchiveTimeRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ArchiveTime, validation.Required, validation.Min(uint32(1)).Error("archive_time must be greater than 0")),
	)
}

func (r UpdatePasswordExpiryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.PasswordExpiry, validation.Required, validation.Min(uint32(1)).Error("password_expiry must be greater than 0")),
	)
}
