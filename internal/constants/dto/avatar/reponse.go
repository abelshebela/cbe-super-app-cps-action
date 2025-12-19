package avatar

import "time"

type AvatarResponseDTO struct {
	ID             string     `json:"id"`
	Avatar         string     `json:"avatar"`
	Label          string     `json:"label"`
	Enable         bool       `json:"enable"`
	IsDeleted      bool       `json:"is_deleted"`
	CreatedAt      time.Time  `json:"created_at"`
	LastModifiedAt time.Time  `json:"last_modified_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
