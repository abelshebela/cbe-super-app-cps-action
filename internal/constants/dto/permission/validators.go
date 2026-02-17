package permission

import (
	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (p CreatePermissionGroupRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.GroupName, validation.Required, validation.Length(1, 50)),
		validation.Field(&p.Role, validation.Required, validation.In("MAKER", "CHECKER")),
		validation.Field(&p.DepartmentID, validation.Required, validation.By(utils.TrimWhiteSpace)),
		validation.Field(&p.PermissionCategoryLists, validation.Required, validation.Each(validation.Required)),
	)
}
func (req UpdatePermissionGroupRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.NewGroupName, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars), validation.Length(1, 50)),
		validation.Field(&req.Role, validation.In("MAKER", "CHECKER")),
		validation.Field(&req.DepartmentID),
	)

}
