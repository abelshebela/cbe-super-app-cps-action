package permission

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)
type CreatePermissionGroupRequest struct {
	GroupName             string   `json:"group_name"`
	Role                  string   `json:"role"`
	PermissionCategoryLists []string `json:"permission_category_list"`
}

func (p CreatePermissionGroupRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.GroupName, validation.Required, validation.Length(1, 50)),
		validation.Field(&p.Role, validation.Required, validation.In("MAKER", "CHECKER")),
		validation.Field(&p.PermissionCategoryLists, validation.Required, validation.Each(validation.Required)),
	)
}