package permission

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (p CreatePermissionGroupRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.GroupName, validation.Required, validation.Length(1, 50)),
		validation.Field(&p.Role, validation.Required, validation.In("MAKER", "CHECKER")),
		validation.Field(&p.PermissionCategoryLists, validation.Required, validation.Each(validation.Required)),
	)
}
func (req *UpdatePermissionGroupRequest) Validate() error {
	if strings.TrimSpace(req.OldGroupName) == "" {
		return fmt.Errorf(localization.ErrorInvalidRequest.Code)
	}
	return nil
}
