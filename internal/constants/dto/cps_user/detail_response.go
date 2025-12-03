package cpsuser

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Detail structs tailored for UI consumption without altering base response structs

type PermissionItemDetail struct {
	ID             bson.ObjectID `json:"_id"`
	PermissionName string        `json:"permissionName"`
}

type PermissionCategoryDetail struct {
	ID           bson.ObjectID          `json:"_id"`
	CategoryName string                 `json:"categoryName"`
	Access       string                 `json:"access"`
	Permissions  []PermissionItemDetail `json:"permissions"`
}

type GroupDetail struct {
	ID                 bson.ObjectID                         `json:"_id"`
	GroupName          string                                `json:"groupName"`
	Permissions        []interface{}                         `json:"permissions"`
	PermissionCategory map[string][]PermissionCategoryDetail `json:"permissionCategory"`
}

type CpsUserDetail struct {
	UserCode           string                                `json:"user_code"`
	FullName           string                                `json:"full_name"`
	Username           string                                `json:"username"`
	UserRole           string                                `json:"role"`
	UserDepartment     string                                `json:"department"`
	UserPortalCards    []string                              `json:"portal_cards"`
	UserPhone          string                                `json:"phone"`
	UserEmail          string                                `json:"email"`
	JoinedAt           time.Time                             `json:"joined_at"`
	UpdatedAt          time.Time                             `json:"updated_at"`
	UserStatus         bool                                  `json:"enabled"`
	UserPermissions    []GroupDetail                         `json:"permissions"`
	PermissionCategory map[string][]PermissionCategoryDetail `json:"permission_category"`
}

// BuildCpsUserDetail builds CpsUserDetail from aggregated CpsUserResponse
func BuildCpsUserDetail(src *CpsUserResponse) *CpsUserDetail {
	if src == nil {
		return nil
	}
	detail := &CpsUserDetail{
		UserCode:           src.UserCode,
		FullName:           src.FullName,
		Username:           src.UserName,
		UserRole:           src.Role,
		UserPhone:          src.PhoneNumber,
		UserEmail:          src.Email,
		JoinedAt:           src.DateJoined,
		UpdatedAt:          src.LastModified,
		UserStatus:         src.Enabled,
		PermissionCategory: map[string][]PermissionCategoryDetail{},
	}
	if src.Department != nil {
		detail.UserDepartment = src.Department.Name
		detail.UserPortalCards = src.Department.PortalCards
	}

	// map permission groups
	groups := make([]GroupDetail, 0, len(src.PermissionGroups))
	for _, g := range src.PermissionGroups {
		gd := GroupDetail{
			ID:                 g.ID,
			GroupName:          g.GroupName,
			Permissions:        []interface{}{},
			PermissionCategory: map[string][]PermissionCategoryDetail{},
		}

		for _, c := range g.PermissionCategory {
			cd := PermissionCategoryDetail{
				ID:           c.ID,
				CategoryName: c.CategoryName,
				Access:       c.Access,
				Permissions:  make([]PermissionItemDetail, 0, len(c.Permissions)),
			}
			for _, p := range c.Permissions {
				cd.Permissions = append(cd.Permissions, PermissionItemDetail{ID: p.ID, PermissionName: p.Name})
			}
			key := strings.ToLower(c.CategoryName)
			gd.PermissionCategory[key] = append(gd.PermissionCategory[key], cd)
			// also accumulate into top-level permission_category
			detail.PermissionCategory[key] = append(detail.PermissionCategory[key], cd)
		}
		groups = append(groups, gd)
	}
	detail.UserPermissions = groups
	return detail
}
