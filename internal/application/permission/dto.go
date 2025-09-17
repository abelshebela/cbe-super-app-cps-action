package permission

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreatePermissionGroupRequest struct {
	GroupName               string   `json:"group_name"`
	Role                    string   `json:"role"`
	PermissionCategoryLists []string `json:"permission_category_list"`
}

func (p CreatePermissionGroupRequest) Validate() error {
	p.Role = strings.ToUpper(p.Role)

	return validation.ValidateStruct(&p,
		validation.Field(&p.GroupName, validation.Required, validation.Length(1, 50)),
		validation.Field(&p.Role, validation.Required, validation.In("MAKER", "CHECKER")),
		validation.Field(&p.PermissionCategoryLists, validation.Required, validation.Each(validation.Required)),
	)
}

type UpdatePermissionGroupRequest struct {
	GroupName               *string   `json:"group_name,omitempty"`
	Role                    *string   `json:"role,omitempty"`
	PermissionCategoryLists *[]string `json:"permission_category_list,omitempty"`
}

func (p UpdatePermissionGroupRequest) Validate() error {
	if p.Role != nil {
		*p.Role = strings.ToUpper(*p.Role)
	}

	return validation.ValidateStruct(&p,
		validation.Field(&p.GroupName, validation.When(p.GroupName != nil, validation.Length(1, 50))),
		validation.Field(&p.Role, validation.When(p.Role != nil, validation.In("MAKER", "CHECKER"))),
		validation.Field(&p.PermissionCategoryLists, validation.When(p.PermissionCategoryLists != nil, validation.Each(validation.Required))),
	)
}
