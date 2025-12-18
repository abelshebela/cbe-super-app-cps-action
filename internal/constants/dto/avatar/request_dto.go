package avatar

import "mime/multipart"

type AvatarDTO struct {
	Label  string                `form:"label" json:"label"`
	Avatar *multipart.FileHeader `form:"avatar" json:"avatar"`
}
