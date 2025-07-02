package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type HQ struct {
	ID        string `bson:"_id" json:"id"`
	Name      string `bson:"name" json:"name"`
	BlockTime uint   `bson:"block_time" json:"block_time"`

	ArchiveTime    uint      `bson:"archive_time" json:"archive_time"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	LastModifiedAt time.Time `bson:"last_modified_at" json:"last_modified_at"`
}

type UpdateBlockTimeRequest struct {
	BlockTime uint `json:"block_time"`
}

func (r UpdateBlockTimeRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BlockTime, validation.Required, validation.Min(uint(1)).Error("block_time must be greater than 0")),
	)
}

type UpdateArchiveTimeRequest struct {
	ArchiveTime uint `json:"archive_time"`
}

func (r UpdateArchiveTimeRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ArchiveTime, validation.Required, validation.Min(uint(1)).Error("archive_time must be greater than 0")),
	)
}
