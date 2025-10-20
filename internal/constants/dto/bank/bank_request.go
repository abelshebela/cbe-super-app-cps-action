package bank_dto

import (
	"mime/multipart"
	"net/http"
)

type CreateBankRequest struct {
	Name string                `form:"name" json:"name"`
	Logo *multipart.FileHeader `form:"logo" json:"logo"`
	Code string                `form:"code" json:"code"`
	BIC  string                `form:"bic" json:"bic"`
}

type UpdateBankRequest struct {
	ID   string                `json:"_id" bson:"_id"`
	Logo *multipart.FileHeader `form:"logo" json:"logo"`
	Name string                `json:"name" bson:"name"`
	Code string                `json:"code" bson:"code"`
	BIC  string                `json:"bic" bson:"bic"`
}

const maxFileSize = 2 * 1024 * 1024 // 2MB

var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

var IsValidImage = func(fileHeader *multipart.FileHeader) bool {
	file, err := fileHeader.Open()
	if err != nil {
		return false
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false
	}
	contentType := http.DetectContentType(buffer)
	return allowedMIMETypes[contentType]
}

type UpdateLogo struct {
	ID   string                `json:"_id" bson:"_id"`
	Logo *multipart.FileHeader `form:"logo" json:"logo"`
}
