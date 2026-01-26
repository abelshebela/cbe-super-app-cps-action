package customer

import (
	"cbe-super-app-cps-action/internal/constants/dto/customer"

	bps "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

func MapToDto(cus *member.User) customer.FindCustomerByIDResponse {
	return customer.FindCustomerByIDResponse{
		ID:                   cus.ID,
		UserCode:             cus.UserCode,
		FullName:             cus.FullName,
		BranchCode:           cus.BranchCode,
		ActivationBranchCode: cus.ActivationBranchCode,
		PhoneNumber:          cus.PhoneNumber,
		Language:             cus.Language,
		Avatar:               cus.Avatar,
		Email:                cus.Email,
		PushToken:            cus.PushToken,
		CustomerNumber:       cus.CustomerNumber,
		UserCategory:         cus.UserCategory,
		Industry:             cus.Industry,
		Sector:               cus.Sector,
		Username:             cus.Username,
		Ownership:            cus.Ownership,
		CustomerSegment:      cus.CustomerSegment,
		BlockedReason:        cus.BlockedReason,
		DeviceUUID:           cus.DeviceUUID,
		AppVersion:           cus.AppVersion,
		Gender:               cus.Gender,
		MemberType:           cus.MemberType,
		Platform:             cus.Platform,
		DeviceStatus:         cus.DeviceStatus,
		OnboardingMethod:     cus.OnboardingMethod,
		KYCLevel:             cus.KYCLevel,
		BlockedOn:            cus.BlockedOn,
		IsUSSDEnabled:        cus.IsUSSDEnabled,
		ISuperappEnabled:     cus.ISuperappEnabled,
		LastLogin:            cus.LastLogin,
		APPInstallationDate:  cus.APPInstallationDate,
		Config:               cus.Config,
		CreatedAt:            cus.CreatedAt,
		ExpiryAt:             cus.ExpiryAt,
		LastModifiedAt:       cus.LastModifiedAt,
		IsBlocked:            cus.IsBlocked,
		IsLocked:             cus.IsLocked,
		Enabled:              cus.Enabled,
		FirstPinSet:          cus.FirstPinSet,
		IsActivated:          cus.IsActivated,
	}
}
func MapBpsActionToCustomerLog(action []bps.BPSAction) []customer.CustomerActionLogResponse {
	response := make([]customer.CustomerActionLogResponse, 0)
	for _, action := range action {
		res := customer.CustomerActionLogResponse{
			ActionCode:     action.ActionCode,
			MakerName:      action.MakerName,
			ActionReason:   action.ActionReason,
			RequestAction:  string(action.RequestAction),
			ServiceName:    action.ServiceName,
			Status:         string(action.Status),
			CreatedAt:      action.CreatedAt,
			LastModifiedAt: action.LastModifiedAt,
		}
		response = append(response, res)
	}
	return response
}
