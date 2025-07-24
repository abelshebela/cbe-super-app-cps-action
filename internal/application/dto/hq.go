package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type HQ struct {
	ID        string `bson:"_id" json:"id"`
	Name      string `bson:"name" json:"name"`
	BlockTime uint32 `bson:"block_time" json:"block_time"`

	ArchiveTime    uint32    `bson:"archive_time" json:"archive_time"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	LastModifiedAt time.Time `bson:"last_modified_at" json:"last_modified_at"`
}

type BlockTimeResponse struct {
	BlockTime      uint32    `json:"block_time"`
	CreatedAtBlock time.Time `json:"created_at_block"`
	UpdatedAtBlock time.Time `json:"updated_at_block"`
}

type ArchiveTimeResponse struct {
	ArchiveTime      uint32    `json:"archive_time"`
	CreatedAtArchive time.Time `json:"created_at_archive"`
	UpdatedAtArchive time.Time `json:"updated_at_archive"`
}

type PasswordExpiryResponse struct {
	PasswordExpiry          uint32    `json:"password_expiry"`
	CreatedAtPasswordExpiry time.Time `json:"created_at_password_expiry"`
	UpdatedAtPasswordExpiry time.Time `json:"updated_at_password_expiry"`
}

type UpdateBlockTimeRequest struct {
	BlockTime uint32 `json:"block_time"`
}

func (r UpdateBlockTimeRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BlockTime, validation.Required, validation.Min(uint32(1)).Error("block_time must be greater than 0")),
	)
}

type UpdateArchiveTimeRequest struct {
	ArchiveTime uint32 `json:"archive_time"`
}

func (r UpdateArchiveTimeRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ArchiveTime, validation.Required, validation.Min(uint32(1)).Error("archive_time must be greater than 0")),
	)
}

type UpdatePasswordExpiryRequest struct {
	PasswordExpiry uint32 `json:"password_expiry"`
}

func (r UpdatePasswordExpiryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.PasswordExpiry, validation.Required, validation.Min(uint32(1)).Error("password_expiry must be greater than 0")),
	)
}
