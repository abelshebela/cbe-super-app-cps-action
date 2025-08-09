package avatar

import (
	"mime/multipart"
	"net/http"
	"time"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AvatarDTO struct {
	Label  string                `form:"label" json:"label"`
	Avatar *multipart.FileHeader `form:"avatar" json:"avatar"`
}

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

func (dto AvatarDTO) IsEmpty() bool {
	return dto.Label == "" && dto.Avatar == nil
}

// MaxAvatarSize defines the max allowed size for avatar: 2MB 
const MaxAvatarSize = 2 * 1024 * 1024

// isImageFormat checks if the content type is an allowed image
func isImageFormat(fileHeader *multipart.FileHeader) bool {
	if fileHeader == nil {
		return false
	}
	contentType := fileHeader.Header.Get("Content-Type")
	switch contentType {
	case "image/jpeg", "image/png", "image/gif":
		return true
	default:
		return false
	}
}

func (dto AvatarDTO) Validate(isCreate bool) error {
	if !isCreate && dto.IsEmpty() {
		return nil
	}

	var rules []*validation.FieldRules

	if isCreate {
		rules = []*validation.FieldRules{
			validation.Field(&dto.Label, validation.Required.Error("label is required")),
			validation.Field(&dto.Avatar,
				validation.Required.Error("avatar image is required"),
				validation.By(validateAvatarFile),
			),
		}
	} else {
		if dto.Label != "" {
			rules = append(rules,
				validation.Field(&dto.Label, validation.Required.Error("label is required")),
			)
		}
		if dto.Avatar != nil {
			rules = append(rules,
				validation.Field(&dto.Avatar, validation.By(validateAvatarFile)),
			)
		}
	}

	if len(rules) > 0 {
		if err := validation.ValidateStruct(&dto, rules...); err != nil {
			return err
		}
	}

	return nil
}

// validateAvatarFile checks the size and format of the uploaded image
func validateAvatarFile(value any) error {
	fileHeader, ok := value.(*multipart.FileHeader)
	if !ok || fileHeader == nil {
		return nil // Nothing to validate
	}

	if fileHeader.Size > MaxAvatarSize {
		return validation.NewError("avatar_file_too_large", "avatar must not exceed 2MB")
	}

	if !isImageFormat(fileHeader) {
		return validation.NewError("invalid_avatar_format", "avatar must be a valid image (jpeg, png, gif)")
	}

	return nil
}

func ParseAvatarRequestFromMultipartForm(r *http.Request, isCreate bool) (AvatarDTO, error) {
	var req AvatarDTO
	file, fileHeader, err := common_util.ParseMultipartFormFile(r, "avatar", 10<<20)
	if err != nil && (isCreate || err.Error() != common_util.ErrMissingFile) {
		return AvatarDTO{}, err
	}
	if file != nil {
		defer file.Close()
		req.Avatar = fileHeader
	}

	req.Label = r.FormValue("label")
	return req, nil
}
