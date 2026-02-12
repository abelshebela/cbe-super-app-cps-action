package topupDto

import (
	"mime/multipart"
)

type TopupRequest struct {
	Name   string                `form:"name" json:"name"`
	Avatar *multipart.FileHeader `form:"avatar" json:"avatar"`
	Code   string                `form:"code" json:"code"`
	// Self   bool                  `form:"self" json:"self"`
	// Other  bool                  `form:"other" json:"other"`
	// Agent  bool                  `form:"agent" json:"agent"`
}
