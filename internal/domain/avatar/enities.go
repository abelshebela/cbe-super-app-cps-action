package avatar

import "time"

type Avatar struct {
	ID             string    `json:"id" bson:"id"`
	AvatarURL      string    `json:"avatar" bson:"avatar"`
	Label          string    `json:"label" bson:"label"`
	Enable         bool      `json:"enable" bson:"enable"`
	IsDeleted      bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      time.Time `json:"deleted_at" bson:"deleted_at"`
}
