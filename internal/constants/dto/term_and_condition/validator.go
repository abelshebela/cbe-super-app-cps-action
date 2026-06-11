package term_and_condition_dto

import (
	"mime"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	tacSafeStringRe = regexp.MustCompile(`^[a-zA-Z0-9 _\-\.]+$`)
)

var allowedDocumentMIMEs = map[string]bool{
	"application/pdf": true,
}

func tacNoSpecialChars(value interface{}) error {
	s, _ := value.(string)
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if !tacSafeStringRe.MatchString(s) {
		return validation.NewError("validation_special_chars", "contains invalid characters")
	}
	return nil
}

func IsValidDocument(req *CreateTACRequest) error {
	if req.TermAndCondition == nil {
		return validation.NewError("validation_file_required", "term_and_condition file is required")
	}

	contentTypeHeader := req.TermAndCondition.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentTypeHeader)
	if err != nil || !allowedDocumentMIMEs[strings.ToLower(mediaType)] {
		return validation.NewError("validation_file_type", "only PDF files are accepted")
	}

	ext := strings.ToLower(req.TermAndCondition.Filename)
	if !strings.HasSuffix(ext, ".pdf") {
		return validation.NewError("validation_file_ext", "file must have a .pdf extension")
	}

	return nil
}

func (r CreateTACRequest) Validate() error {
	if err := IsValidDocument(&r); err != nil {
		return err
	}
	return validation.ValidateStruct(&r,
		validation.Field(&r.AccountProductID,
			validation.Required,
		),
		validation.Field(&r.ActivationTime,
			validation.Required,
			validation.Length(1, 128),
			validation.By(tacNoSpecialChars),
		),
		validation.Field(&r.VersionLabel,
			validation.Required,
			validation.Length(1, 128),
			validation.By(tacNoSpecialChars),
		),
	)
}
