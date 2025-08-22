package ad

import (
	"cbe-super-app-cps-action/internal/constants/model"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	"context"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)


func BuildCPSAction(ctx context.Context, request cpsaction.CreateCPSRequest) *model.CPSAction {
	return &model.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          request.User.UserCode,
		MakerName:        request.User.FullName,
		MakerPhoneNumber: request.User.PhoneNumber,
		Department:       request.User.Department,
		CurrentAction:    request.CurData,
		PreviousAction:   request.PrevData,
		RequestAction:    string(request.RequestAction),
		ActionStatus:     string(request.ActionStatus),
		ActionType:       string(request.ActionType),
		MakerActionTime:  time.Now(),
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}
}


// generateAdvert creates a copy of an advert
func GenerateAdvert(advert model.Advert) *model.Advert {
	return &model.Advert{
		ID:            advert.ID,
		Title:         advert.Title,
		Description:   advert.Description,
		BannerImage:   advert.BannerImage,
		AdvertFor:     advert.AdvertFor,
		Date:          advert.Date,
		Enabled:       advert.Enabled,
		IsDeleted:     advert.IsDeleted,
		CreatedAt:     advert.CreatedAt,
		LastUpdatedAt: advert.LastUpdatedAt,
		DeletedAt:     advert.DeletedAt,
	}
}
