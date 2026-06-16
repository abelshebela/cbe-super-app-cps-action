package term_and_condition_dto

import (
	"mime"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	tacSafeStringRe = regexp.MustCompile(`^[a-zA-Z0-9 _\-\.]+$`)
)

var allowedDocumentMIMEs = map[string]bool{
	"application/pdf": true,
}

func tacValidDate(value interface{}) error {
	s, _ := value.(string)
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02",
		"Mon Jan 02 2006 15:04:05 MST-0700 (MST)",
	}
	for _, f := range formats {
		if _, err := time.Parse(f, s); err == nil {
			return nil
		}
	}

	return validation.NewError("validation_date_format", "must be a valid RFC3339 timestamp (e.g., '2026-06-15T00:00:00Z') or date (e.g., '2026-06-15')")
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

func (r UpdateTACRequest) Validate() error {
	if r.TermAndCondition != nil {
		if err := IsValidDocument(&CreateTACRequest{TermAndCondition: r.TermAndCondition}); err != nil {
			return err
		}
	}
	return validation.ValidateStruct(&r,
		validation.Field(&r.ActivationTime,
			validation.By(tacValidDate),
		),
		validation.Field(&r.VersionLabel,
			validation.Length(0, 128),
			validation.By(tacNoSpecialChars),
		),
	)
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
			validation.By(tacValidDate),
		),
		validation.Field(&r.VersionLabel,
			validation.Required,
			validation.Length(1, 128),
			validation.By(tacNoSpecialChars),
		),
	)
}
