package bps_actionrole_dto

import (
	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r CreateActionRoleRequest) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.ActionName,
			validation.Required.Error("Action name is required"),
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.PortalCardName,
			validation.Required.Error("Portal card name is required"),
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
			validation.Length(0, 50).Error("Portal card name must not exceed 50 characters"),
		),
	)
	if err != nil {
		return err
	}

	if r.IsViewOnly {
		return nil
	}

	if err := validation.Validate(r.AssignedMakersRoles,
		validation.Required.Error("Assigned maker role is required"),
		validation.Length(1, 0).Error("Assigned maker role must have at least 1 role"),
		validation.Each(
			validation.Required,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),
	); err != nil {
		return err
	}

	if err := validate2DStringSliceRequired(r.AssignedAuditorRoles, "assigned_auditor_roles"); err != nil {
		return err
	}

	if r.IsMakerOnly {
		return nil
	}

	if err := validate2DStringSliceRequired(r.AssignedCheckerRoles, "assigned_checkers_roles"); err != nil {
		return err
	}

	return nil
}

func (r UpdateActionRoleRequest) Validate() error {

	return validation.ValidateStruct(&r,
		validation.Field(&r.AssignedMakersRoles,
			validation.Each(
				validation.Required.Error("Assigned maker role item cannot be empty"),
				validation.By(utils.TrimWhiteSpace),
				validation.By(utils.NoSpecialChars),
			),
		),

		validation.Field(&r.AssignedCheckerRoles,
			validation.Each(
				validation.Required,
				validation.Each(
					validation.Required,
					validation.By(utils.TrimWhiteSpace),
					validation.By(utils.NoSpecialChars),
				),
			),
		),

		validation.Field(&r.AssignedAuditorRoles,
			validation.Each(
				validation.Required,
				validation.Each(
					validation.Required,
					validation.By(utils.TrimWhiteSpace),
					validation.By(utils.NoSpecialChars),
				),
			),
		),
	)
}

func validate2DStringSliceRequired(value [][]string, field string) error {
	if err := validation.Validate(value,
		validation.Required.Error(field+" is required"),
		validation.Length(1, 0).Error(field+" must have at least 1 group"),
		validation.Each(
			validation.Required.Error(field+" group cannot be empty"),
			validation.Length(1, 0).Error(field+" group must have at least 1 role"),
			validation.Each(
				validation.Required.Error(field+" item cannot be empty"),
				validation.By(utils.TrimWhiteSpace),
				validation.By(utils.NoSpecialChars),
			),
		),
	); err != nil {
		return err
	}
	return nil
}
