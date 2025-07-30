// Package dto provides data transfer objects for wallet operations.
package wallet

import (
	"mime/multipart"
)

type WalletRequest struct {
	Name   string                `form:"name" json:"name"`
	Avatar *multipart.FileHeader `form:"avatar" json:"avatar"`
	Code   string                `form:"code" json:"code"`
}
