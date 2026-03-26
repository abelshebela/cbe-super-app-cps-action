package notification

import (
	"regexp"
	"strings"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/text/unicode/norm"
)

// var safePattern = regexp.MustCompile(`^[\p{L}\p{M}\p{N}\s._@'-]+$`)
var safePattern = regexp.MustCompile(
	`^[\p{L}\p{N}\p{P}\p{Zs}\n\r\t]+$`,
)

func noDangerousChars(value any) error {
	var str string

	switch v := value.(type) {
	case string:
		str = v
	case *string:
		if v == nil {
			return nil
		}
		str = *v
	default:
		return validation.NewError("validation", "invalid type")
	}

	// Normalize Unicode (blocks homoglyph / confusable tricks)
	str = norm.NFKC.String(str)

	str = strings.TrimSpace(str)
	if str == "" {
		return nil
	}

	// Allow letters, numbers, punctuation, spaces, and newlines
	if !safePattern.MatchString(str) {
		return validation.NewError("validation", "contains invalid characters")
	}

	// Explicitly block hidden / control characters (except whitespace)
	for _, r := range str {
		if unicode.IsControl(r) && !unicode.IsSpace(r) {
			return validation.NewError("validation", "contains control characters")
		}
	}

	return nil
}

var allowedNotificationFor = []string{"IFB", "CB", "ALL"}

func (r NotificationRequest) Validate(isCreate bool) error {
	enumValues := toInterfaceSlice(allowedNotificationFor)

	var fieldRules []*validation.FieldRules

	if isCreate {
		fieldRules = []*validation.FieldRules{
			validation.Field(&r.NotificationType, validation.Required.Error("notification_type is required"), validation.By(noDangerousChars)),
			validation.Field(&r.NotificationBody,
				validation.Required.Error("notification_body is required"),
				validation.By(noDangerousChars),
				validation.Length(1, 200),
			),
			validation.Field(&r.Title, validation.Required.Error("title is required"), validation.By(noDangerousChars)),
			validation.Field(&r.For,
				validation.Required.Error("for is required"),
				validation.By(noDangerousChars),
				validation.In(enumValues...).Error("invalid value for 'for'"),
			),
		}
	} else {
		fieldRules = []*validation.FieldRules{
			validation.Field(&r.NotificationType, validation.By(noDangerousChars)),
			validation.Field(&r.NotificationBody, validation.By(noDangerousChars)),
			validation.Field(&r.Title, validation.By(noDangerousChars)),
			validation.Field(&r.For,
				validation.By(noDangerousChars),
				validation.In(enumValues...).Error("invalid value for 'for'"),
			),
		}
	}

	return validation.ValidateStruct(&r, fieldRules...)
}

func toInterfaceSlice(strs []string) []interface{} {
	res := make([]interface{}, len(strs))
	for i, s := range strs {
		res[i] = s
	}
	return res
}

func (r NotificationRequest) IsEmpty() bool {
	return strings.TrimSpace(r.NotificationType) == "" &&
		strings.TrimSpace(r.NotificationBody) == "" &&
		strings.TrimSpace(r.For) == "" &&
		strings.TrimSpace(r.Title) == ""
}
