package avatar

import (
	"mime/multipart"
)

type AvatarRequest struct {
	Label  string                `form:"label"`
	Avatar *multipart.FileHeader `form:"avatar"`
}
